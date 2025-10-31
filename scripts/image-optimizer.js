#!/usr/bin/env node

/**
 * Image Optimization System
 * Converts images to WebP format and ensures all images are under 200KB
 * Supports batch processing, quality optimization, and validation
 */

const fs = require('fs').promises;
const path = require('path');
const { execSync, spawn } = require('child_process');
const crypto = require('crypto');

// Configuration
const CONFIG = {
  maxFileSize: 200 * 1024, // 200KB in bytes
  targetFormats: ['webp'],
  sourceFormats: ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'tiff'],
  qualityLevels: [90, 80, 70, 60, 50, 40, 30],
  directories: {
    source: './frontend/public/images',
    output: './frontend/public/images/optimized',
    backup: './frontend/public/images/backup',
    temp: './temp/image-processing'
  },
  webpOptions: {
    quality: 80,
    method: 6, // Compression method (0-6, 6 is slowest but best)
    lossless: false,
    nearLossless: false,
    alphaQuality: 100,
    autoFilter: true,
    sharpness: 0,
    filterStrength: 60,
    filterSharpness: 0,
    filterType: 1,
    partitions: 0,
    segments: 4,
    pass: 1,
    showCompressed: false,
    preprocessing: 0,
    partitionLimit: 0,
    alphaMethod: 1,
    alphaFilter: 'fast',
    exact: false,
    blend: true,
    noAlpha: false,
    hint: 'default'
  },
  resizeOptions: {
    maxWidth: 1920,
    maxHeight: 1080,
    thumbnailSizes: [150, 300, 600, 1200]
  },
  validation: {
    checkDuplicates: true,
    generateHashes: true,
    createManifest: true,
    validateSizes: true,
    checkQuality: true
  }
};

// Image Optimizer Class
class ImageOptimizer {
  constructor(config = CONFIG) {
    this.config = config;
    this.stats = {
      processed: 0,
      optimized: 0,
      failed: 0,
      totalSizeBefore: 0,
      totalSizeAfter: 0,
      duplicatesFound: 0,
      errors: [],
      warnings: [],
      processingTime: 0
    };
    this.manifest = {
      version: '1.0.0',
      generatedAt: new Date().toISOString(),
      images: [],
      duplicates: [],
      statistics: {}
    };
    this.imageHashes = new Map();
  }

  /**
   * Main optimization process
   */
  async optimize(options = {}) {
    const startTime = Date.now();
    
    try {
      console.log('🖼️  Starting Image Optimization Process...');
      
      // Setup directories
      await this.setupDirectories();
      
      // Check dependencies
      await this.checkDependencies();
      
      // Find all images
      const images = await this.findImages();
      console.log(`📁 Found ${images.length} images to process`);
      
      if (images.length === 0) {
        console.log('ℹ️  No images found to optimize');
        return this.getReport();
      }
      
      // Process images
      await this.processImages(images, options);
      
      // Validate results
      await this.validateResults();
      
      // Generate manifest
      await this.generateManifest();
      
      // Cleanup
      await this.cleanup();
      
      this.stats.processingTime = Date.now() - startTime;
      
      console.log('✅ Image optimization completed!');
      return this.getReport();
      
    } catch (error) {
      this.stats.errors.push(`Optimization failed: ${error.message}`);
      console.error('❌ Optimization failed:', error);
      throw error;
    }
  }

  /**
   * Setup required directories
   */
  async setupDirectories() {
    const dirs = Object.values(this.config.directories);
    
    for (const dir of dirs) {
      try {
        await fs.mkdir(dir, { recursive: true });
      } catch (error) {
        if (error.code !== 'EEXIST') {
          throw new Error(`Failed to create directory ${dir}: ${error.message}`);
        }
      }
    }
  }

  /**
   * Check required dependencies
   */
  async checkDependencies() {
    const dependencies = ['cwebp', 'identify'];
    const missing = [];
    
    for (const dep of dependencies) {
      try {
        execSync(`which ${dep}`, { stdio: 'ignore' });
      } catch (error) {
        missing.push(dep);
      }
    }
    
    if (missing.length > 0) {
      const installCmd = process.platform === 'darwin' 
        ? `brew install webp imagemagick`
        : `sudo apt-get install webp imagemagick`;
      
      throw new Error(
        `Missing dependencies: ${missing.join(', ')}\n` +
        `Install with: ${installCmd}`
      );
    }
  }

  /**
   * Find all images in source directories
   */
  async findImages() {
    const images = [];
    
    const findInDirectory = async (dir) => {
      try {
        const entries = await fs.readdir(dir, { withFileTypes: true });
        
        for (const entry of entries) {
          const fullPath = path.join(dir, entry.name);
          
          if (entry.isDirectory()) {
            // Skip output and backup directories
            if (!fullPath.includes('optimized') && !fullPath.includes('backup')) {
              const subImages = await findInDirectory(fullPath);
              images.push(...subImages);
            }
          } else if (entry.isFile()) {
            const ext = path.extname(entry.name).toLowerCase().slice(1);
            if (this.config.sourceFormats.includes(ext)) {
              const stats = await fs.stat(fullPath);
              images.push({
                path: fullPath,
                name: entry.name,
                extension: ext,
                size: stats.size,
                directory: dir,
                relativePath: path.relative(this.config.directories.source, fullPath)
              });
            }
          }
        }
      } catch (error) {
        if (error.code !== 'ENOENT') {
          this.stats.warnings.push(`Failed to read directory ${dir}: ${error.message}`);
        }
      }
    };
    
    await findInDirectory(this.config.directories.source);
    return images;
  }

  /**
   * Process all images
   */
  async processImages(images, options = {}) {
    const concurrency = options.concurrency || 4;
    const batches = this.createBatches(images, concurrency);
    
    for (let i = 0; i < batches.length; i++) {
      const batch = batches[i];
      console.log(`📦 Processing batch ${i + 1}/${batches.length} (${batch.length} images)`);
      
      const promises = batch.map(image => this.processImage(image, options));
      await Promise.all(promises);
    }
  }

  /**
   * Create batches for concurrent processing
   */
  createBatches(items, batchSize) {
    const batches = [];
    for (let i = 0; i < items.length; i += batchSize) {
      batches.push(items.slice(i, i + batchSize));
    }
    return batches;
  }

  /**
   * Process a single image
   */
  async processImage(image, options = {}) {
    try {
      console.log(`🔄 Processing: ${image.name}`);
      
      this.stats.processed++;
      this.stats.totalSizeBefore += image.size;
      
      // Generate hash for duplicate detection
      const hash = await this.generateImageHash(image.path);
      
      // Check for duplicates
      if (this.config.validation.checkDuplicates && this.imageHashes.has(hash)) {
        const duplicate = this.imageHashes.get(hash);
        this.stats.duplicatesFound++;
        this.manifest.duplicates.push({
          original: duplicate.path,
          duplicate: image.path,
          hash: hash
        });
        console.log(`🔍 Duplicate found: ${image.name} (matches ${duplicate.name})`);
        return;
      }
      
      this.imageHashes.set(hash, image);
      
      // Create backup if requested
      if (options.createBackup) {
        await this.createBackup(image);
      }
      
      // Get image metadata
      const metadata = await this.getImageMetadata(image.path);
      
      // Resize if needed
      let processedPath = image.path;
      if (this.shouldResize(metadata)) {
        processedPath = await this.resizeImage(image, metadata);
      }
      
      // Convert to WebP with quality optimization
      const optimizedPath = await this.convertToWebP(processedPath, image);
      
      // Validate result
      const optimizedStats = await fs.stat(optimizedPath);
      this.stats.totalSizeAfter += optimizedStats.size;
      
      if (optimizedStats.size <= this.config.maxFileSize) {
        this.stats.optimized++;
        
        // Add to manifest
        this.manifest.images.push({
          original: {
            path: image.relativePath,
            name: image.name,
            size: image.size,
            format: image.extension
          },
          optimized: {
            path: path.relative(this.config.directories.source, optimizedPath),
            size: optimizedStats.size,
            format: 'webp',
            compression: ((image.size - optimizedStats.size) / image.size * 100).toFixed(2) + '%'
          },
          metadata: metadata,
          hash: hash,
          processedAt: new Date().toISOString()
        });
        
        console.log(`✅ Optimized: ${image.name} (${this.formatBytes(image.size)} → ${this.formatBytes(optimizedStats.size)})`);
      } else {
        this.stats.warnings.push(`Image ${image.name} still exceeds size limit after optimization`);
        console.log(`⚠️  Warning: ${image.name} still too large (${this.formatBytes(optimizedStats.size)})`);
      }
      
      // Generate thumbnails if requested
      if (options.generateThumbnails) {
        await this.generateThumbnails(optimizedPath, image);
      }
      
    } catch (error) {
      this.stats.failed++;
      this.stats.errors.push(`Failed to process ${image.name}: ${error.message}`);
      console.error(`❌ Failed to process ${image.name}:`, error.message);
    }
  }

  /**
   * Generate image hash for duplicate detection
   */
  async generateImageHash(imagePath) {
    const buffer = await fs.readFile(imagePath);
    return crypto.createHash('md5').update(buffer).digest('hex');
  }

  /**
   * Create backup of original image
   */
  async createBackup(image) {
    const backupPath = path.join(
      this.config.directories.backup,
      image.relativePath
    );
    
    const backupDir = path.dirname(backupPath);
    await fs.mkdir(backupDir, { recursive: true });
    
    await fs.copyFile(image.path, backupPath);
  }

  /**
   * Get image metadata using ImageMagick
   */
  async getImageMetadata(imagePath) {
    try {
      const output = execSync(
        `identify -format "%w,%h,%[colorspace],%Q,%[compression]" "${imagePath}"`,
        { encoding: 'utf8' }
      ).trim();
      
      const [width, height, colorspace, quality, compression] = output.split(',');
      
      return {
        width: parseInt(width),
        height: parseInt(height),
        colorspace: colorspace || 'Unknown',
        quality: parseInt(quality) || 0,
        compression: compression || 'Unknown'
      };
    } catch (error) {
      console.warn(`Failed to get metadata for ${imagePath}:`, error.message);
      return {
        width: 0,
        height: 0,
        colorspace: 'Unknown',
        quality: 0,
        compression: 'Unknown'
      };
    }
  }

  /**
   * Check if image should be resized
   */
  shouldResize(metadata) {
    return metadata.width > this.config.resizeOptions.maxWidth ||
           metadata.height > this.config.resizeOptions.maxHeight;
  }

  /**
   * Resize image if needed
   */
  async resizeImage(image, metadata) {
    const tempPath = path.join(
      this.config.directories.temp,
      `resized_${Date.now()}_${image.name}`
    );
    
    const { maxWidth, maxHeight } = this.config.resizeOptions;
    
    try {
      execSync(
        `convert "${image.path}" -resize ${maxWidth}x${maxHeight}> "${tempPath}"`,
        { stdio: 'ignore' }
      );
      
      console.log(`📏 Resized: ${image.name} (${metadata.width}x${metadata.height} → ${maxWidth}x${maxHeight})`);
      return tempPath;
    } catch (error) {
      throw new Error(`Failed to resize image: ${error.message}`);
    }
  }

  /**
   * Convert image to WebP format with quality optimization
   */
  async convertToWebP(imagePath, originalImage) {
    const outputPath = path.join(
      this.config.directories.output,
      originalImage.relativePath.replace(/\.[^.]+$/, '.webp')
    );
    
    const outputDir = path.dirname(outputPath);
    await fs.mkdir(outputDir, { recursive: true });
    
    // Try different quality levels to meet size requirement
    for (const quality of this.config.qualityLevels) {
      const tempPath = `${outputPath}.tmp`;
      
      try {
        // Build cwebp command
        const options = {
          ...this.config.webpOptions,
          quality: quality
        };
        
        const cmd = this.buildWebPCommand(imagePath, tempPath, options);
        execSync(cmd, { stdio: 'ignore' });
        
        // Check file size
        const stats = await fs.stat(tempPath);
        
        if (stats.size <= this.config.maxFileSize) {
          // Size is acceptable, use this version
          await fs.rename(tempPath, outputPath);
          console.log(`🎯 WebP conversion successful at quality ${quality}`);
          return outputPath;
        } else {
          // Size too large, try next quality level
          await fs.unlink(tempPath);
          console.log(`📊 Quality ${quality} too large (${this.formatBytes(stats.size)}), trying lower quality`);
        }
      } catch (error) {
        // Clean up temp file if it exists
        try {
          await fs.unlink(tempPath);
        } catch (e) {}
        
        if (quality === this.config.qualityLevels[this.config.qualityLevels.length - 1]) {
          throw new Error(`WebP conversion failed at all quality levels: ${error.message}`);
        }
      }
    }
    
    throw new Error('Could not optimize image to meet size requirements');
  }

  /**
   * Build cwebp command with options
   */
  buildWebPCommand(inputPath, outputPath, options) {
    let cmd = `cwebp "${inputPath}" -o "${outputPath}"`;
    
    // Add quality
    cmd += ` -q ${options.quality}`;
    
    // Add method
    cmd += ` -m ${options.method}`;
    
    // Add other options
    if (options.lossless) cmd += ' -lossless';
    if (options.nearLossless) cmd += ` -near_lossless ${options.nearLossless}`;
    if (options.alphaQuality !== 100) cmd += ` -alpha_q ${options.alphaQuality}`;
    if (options.autoFilter) cmd += ' -af';
    if (options.sharpness !== 0) cmd += ` -sharpness ${options.sharpness}`;
    if (options.filterStrength !== 60) cmd += ` -f ${options.filterStrength}`;
    if (options.noAlpha) cmd += ' -noalpha';
    
    return cmd;
  }

  /**
   * Generate thumbnails for responsive images
   */
  async generateThumbnails(imagePath, originalImage) {
    const thumbnailDir = path.join(
      this.config.directories.output,
      'thumbnails',
      path.dirname(originalImage.relativePath)
    );
    
    await fs.mkdir(thumbnailDir, { recursive: true });
    
    const baseName = path.basename(originalImage.name, path.extname(originalImage.name));
    
    for (const size of this.config.resizeOptions.thumbnailSizes) {
      const thumbnailPath = path.join(thumbnailDir, `${baseName}_${size}w.webp`);
      
      try {
        const cmd = `cwebp "${imagePath}" -o "${thumbnailPath}" -resize ${size} 0 -q 80`;
        execSync(cmd, { stdio: 'ignore' });
        
        console.log(`🖼️  Generated thumbnail: ${size}w`);
      } catch (error) {
        this.stats.warnings.push(`Failed to generate ${size}w thumbnail for ${originalImage.name}`);
      }
    }
  }

  /**
   * Validate optimization results
   */
  async validateResults() {
    console.log('🔍 Validating optimization results...');
    
    const outputImages = await this.findImages();
    let oversizedCount = 0;
    
    for (const image of outputImages) {
      if (image.size > this.config.maxFileSize) {
        oversizedCount++;
        this.stats.warnings.push(
          `Image ${image.name} exceeds size limit: ${this.formatBytes(image.size)}`
        );
      }
    }
    
    if (oversizedCount > 0) {
      console.log(`⚠️  ${oversizedCount} images still exceed size limit`);
    } else {
      console.log('✅ All images meet size requirements');
    }
  }

  /**
   * Generate optimization manifest
   */
  async generateManifest() {
    this.manifest.statistics = {
      totalProcessed: this.stats.processed,
      totalOptimized: this.stats.optimized,
      totalFailed: this.stats.failed,
      duplicatesFound: this.stats.duplicatesFound,
      sizeBefore: this.stats.totalSizeBefore,
      sizeAfter: this.stats.totalSizeAfter,
      compressionRatio: ((this.stats.totalSizeBefore - this.stats.totalSizeAfter) / this.stats.totalSizeBefore * 100).toFixed(2) + '%',
      processingTime: this.stats.processingTime,
      errors: this.stats.errors,
      warnings: this.stats.warnings
    };
    
    const manifestPath = path.join(this.config.directories.output, 'optimization-manifest.json');
    await fs.writeFile(manifestPath, JSON.stringify(this.manifest, null, 2));
    
    console.log(`📋 Manifest generated: ${manifestPath}`);
  }

  /**
   * Cleanup temporary files
   */
  async cleanup() {
    try {
      const tempFiles = await fs.readdir(this.config.directories.temp);
      for (const file of tempFiles) {
        await fs.unlink(path.join(this.config.directories.temp, file));
      }
    } catch (error) {
      // Ignore cleanup errors
    }
  }

  /**
   * Get optimization report
   */
  getReport() {
    const compressionRatio = this.stats.totalSizeBefore > 0 
      ? ((this.stats.totalSizeBefore - this.stats.totalSizeAfter) / this.stats.totalSizeBefore * 100)
      : 0;
    
    return {
      summary: {
        processed: this.stats.processed,
        optimized: this.stats.optimized,
        failed: this.stats.failed,
        duplicates: this.stats.duplicatesFound,
        processingTime: this.formatTime(this.stats.processingTime)
      },
      sizes: {
        before: this.formatBytes(this.stats.totalSizeBefore),
        after: this.formatBytes(this.stats.totalSizeAfter),
        saved: this.formatBytes(this.stats.totalSizeBefore - this.stats.totalSizeAfter),
        compressionRatio: compressionRatio.toFixed(2) + '%'
      },
      issues: {
        errors: this.stats.errors,
        warnings: this.stats.warnings
      },
      manifest: this.manifest
    };
  }

  /**
   * Format bytes to human readable format
   */
  formatBytes(bytes) {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  }

  /**
   * Format time to human readable format
   */
  formatTime(ms) {
    if (ms < 1000) return `${ms}ms`;
    if (ms < 60000) return `${(ms / 1000).toFixed(1)}s`;
    return `${(ms / 60000).toFixed(1)}m`;
  }
}

// Validation Script
class ImageValidator {
  constructor(config = CONFIG) {
    this.config = config;
  }

  /**
   * Validate all images meet requirements
   */
  async validate() {
    console.log('🔍 Starting Image Validation...');
    
    const results = {
      totalImages: 0,
      validImages: 0,
      oversizedImages: [],
      unsupportedFormats: [],
      missingImages: [],
      duplicates: [],
      errors: []
    };
    
    try {
      // Find all images
      const optimizer = new ImageOptimizer(this.config);
      const images = await optimizer.findImages();
      
      results.totalImages = images.length;
      
      for (const image of images) {
        // Check file size
        if (image.size > this.config.maxFileSize) {
          results.oversizedImages.push({
            path: image.relativePath,
            size: optimizer.formatBytes(image.size),
            limit: optimizer.formatBytes(this.config.maxFileSize)
          });
        } else {
          results.validImages++;
        }
        
        // Check format
        if (!this.config.targetFormats.includes(image.extension) && 
            !this.config.sourceFormats.includes(image.extension)) {
          results.unsupportedFormats.push({
            path: image.relativePath,
            format: image.extension
          });
        }
      }
      
      // Generate report
      const report = this.generateValidationReport(results);
      
      console.log('✅ Image validation completed');
      return report;
      
    } catch (error) {
      results.errors.push(`Validation failed: ${error.message}`);
      console.error('❌ Validation failed:', error);
      return results;
    }
  }

  /**
   * Generate validation report
   */
  generateValidationReport(results) {
    const passRate = results.totalImages > 0 
      ? (results.validImages / results.totalImages * 100).toFixed(1)
      : 0;
    
    const status = results.oversizedImages.length === 0 ? 'PASS' : 'FAIL';
    
    return {
      status,
      summary: {
        total: results.totalImages,
        valid: results.validImages,
        oversized: results.oversizedImages.length,
        unsupported: results.unsupportedFormats.length,
        passRate: `${passRate}%`
      },
      issues: {
        oversizedImages: results.oversizedImages,
        unsupportedFormats: results.unsupportedFormats,
        errors: results.errors
      },
      recommendations: this.generateRecommendations(results)
    };
  }

  /**
   * Generate recommendations based on validation results
   */
  generateRecommendations(results) {
    const recommendations = [];
    
    if (results.oversizedImages.length > 0) {
      recommendations.push(
        `Optimize ${results.oversizedImages.length} oversized images using the image optimizer`
      );
    }
    
    if (results.unsupportedFormats.length > 0) {
      recommendations.push(
        `Convert ${results.unsupportedFormats.length} images to supported formats (WebP recommended)`
      );
    }
    
    if (results.totalImages === 0) {
      recommendations.push('No images found. Check the source directory configuration.');
    }
    
    if (recommendations.length === 0) {
      recommendations.push('All images meet the requirements. Great job!');
    }
    
    return recommendations;
  }
}

// CLI Interface
if (require.main === module) {
  const args = process.argv.slice(2);
  const command = args[0] || 'optimize';
  
  const options = {
    createBackup: args.includes('--backup'),
    generateThumbnails: args.includes('--thumbnails'),
    concurrency: parseInt(args.find(arg => arg.startsWith('--concurrency='))?.split('=')[1]) || 4
  };
  
  async function main() {
    try {
      switch (command) {
        case 'optimize':
          const optimizer = new ImageOptimizer();
          const report = await optimizer.optimize(options);
          
          console.log('\n📊 Optimization Report:');
          console.log(`   Processed: ${report.summary.processed} images`);
          console.log(`   Optimized: ${report.summary.optimized} images`);
          console.log(`   Failed: ${report.summary.failed} images`);
          console.log(`   Size reduction: ${report.sizes.compressionRatio}`);
          console.log(`   Total saved: ${report.sizes.saved}`);
          console.log(`   Processing time: ${report.summary.processingTime}`);
          
          if (report.issues.errors.length > 0) {
            console.log('\n❌ Errors:');
            report.issues.errors.forEach(error => console.log(`   ${error}`));
          }
          
          if (report.issues.warnings.length > 0) {
            console.log('\n⚠️  Warnings:');
            report.issues.warnings.forEach(warning => console.log(`   ${warning}`));
          }
          
          break;
          
        case 'validate':
          const validator = new ImageValidator();
          const validationReport = await validator.validate();
          
          console.log(`\n📋 Validation Report: ${validationReport.status}`);
          console.log(`   Total images: ${validationReport.summary.total}`);
          console.log(`   Valid images: ${validationReport.summary.valid}`);
          console.log(`   Oversized: ${validationReport.summary.oversized}`);
          console.log(`   Pass rate: ${validationReport.summary.passRate}`);
          
          if (validationReport.issues.oversizedImages.length > 0) {
            console.log('\n📏 Oversized Images:');
            validationReport.issues.oversizedImages.forEach(img => {
              console.log(`   ${img.path} (${img.size} > ${img.limit})`);
            });
          }
          
          if (validationReport.recommendations.length > 0) {
            console.log('\n💡 Recommendations:');
            validationReport.recommendations.forEach(rec => console.log(`   ${rec}`));
          }
          
          process.exit(validationReport.status === 'PASS' ? 0 : 1);
          break;
          
        default:
          console.log('Usage:');
          console.log('  node image-optimizer.js optimize [--backup] [--thumbnails] [--concurrency=4]');
          console.log('  node image-optimizer.js validate');
          process.exit(1);
      }
    } catch (error) {
      console.error('❌ Error:', error.message);
      process.exit(1);
    }
  }
  
  main();
}

module.exports = { ImageOptimizer, ImageValidator, CONFIG };