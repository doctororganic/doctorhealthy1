// Halal Filter System for Nutrition Platform
// This module provides comprehensive halal filtering for ingredients and recipes

// Comprehensive list of non-halal ingredients
const NON_HALAL_INGREDIENTS = {
    // Pork and pork-derived products
    pork: [
        'pork', 'bacon', 'ham', 'sausage', 'pepperoni', 'prosciutto', 'pancetta',
        'chorizo', 'salami', 'mortadella', 'lard', 'pork fat', 'pork gelatin',
        'pork chops', 'pork shoulder', 'pork tenderloin', 'pork ribs',
        'خنزير', 'لحم خنزير', 'بيكون', 'هام', 'سجق', 'شحم خنزير'
    ],
    
    // Alcohol and alcoholic beverages
    alcohol: [
        'wine', 'beer', 'whiskey', 'vodka', 'rum', 'gin', 'brandy', 'champagne',
        'sake', 'mirin', 'cooking wine', 'red wine', 'white wine', 'sherry',
        'port wine', 'marsala', 'bourbon', 'scotch', 'tequila', 'liqueur',
        'amaretto', 'kahlua', 'baileys', 'cognac', 'absinthe', 'vermouth',
        'كحول', 'خمر', 'نبيذ', 'بيرة', 'ويسكي', 'فودكا', 'روم', 'جين',
        'براندي', 'شامبانيا', 'ساكي', 'مارسالا', 'شيري'
    ],
    
    // Non-halal gelatin and additives
    gelatin: [
        'gelatin', 'gelatine', 'pork gelatin', 'beef gelatin from non-halal source',
        'carmine', 'cochineal', 'natural red 4', 'e120', 'shellac', 'confectioner\'s glaze',
        'جيلاتين', 'جلاتين خنزير', 'قرمزي', 'صمغ اللك'
    ],
    
    // Non-halal meat and poultry (if not halal certified)
    meat: [
        'non-halal beef', 'non-halal chicken', 'non-halal lamb', 'non-halal turkey',
        'non-halal duck', 'non-halal goose', 'venison', 'rabbit meat',
        'لحم غير حلال', 'دجاج غير حلال', 'لحم غزال', 'لحم أرنب'
    ],
    
    // Enzymes and additives that may be non-halal
    enzymes: [
        'rennet', 'animal rennet', 'pepsin', 'lipase', 'trypsin',
        'pancreatin', 'chymosin', 'animal enzymes',
        'منفحة حيوانية', 'إنزيمات حيوانية'
    ],
    
    // Vanilla extract with alcohol
    extracts: [
        'vanilla extract with alcohol', 'rum extract', 'brandy extract',
        'wine extract', 'beer extract',
        'خلاصة الفانيليا بالكحول', 'خلاصة الروم'
    ]
};

// Halal alternatives for non-halal ingredients
const HALAL_ALTERNATIVES = {
    // Pork alternatives
    'pork': 'halal beef or lamb',
    'bacon': 'halal beef bacon or turkey bacon',
    'ham': 'halal turkey ham or beef ham',
    'sausage': 'halal beef or chicken sausage',
    'pepperoni': 'halal beef pepperoni',
    'lard': 'vegetable oil or ghee',
    'pork gelatin': 'halal beef gelatin or agar-agar',
    
    // Alcohol alternatives
    'wine': 'grape juice or pomegranate juice',
    'red wine': 'red grape juice or cranberry juice',
    'white wine': 'white grape juice or apple juice',
    'beer': 'non-alcoholic malt beverage',
    'cooking wine': 'grape juice with vinegar',
    'mirin': 'rice vinegar with sugar',
    'sake': 'rice vinegar',
    'brandy': 'apple juice concentrate',
    'rum': 'pineapple juice concentrate',
    'whiskey': 'apple cider vinegar',
    
    // Gelatin alternatives
    'gelatin': 'halal gelatin or agar-agar',
    'pork gelatin': 'halal beef gelatin or agar-agar',
    'carmine': 'natural beetroot extract',
    'cochineal': 'natural red food coloring',
    
    // Extract alternatives
    'vanilla extract with alcohol': 'alcohol-free vanilla extract',
    'rum extract': 'pineapple flavoring',
    'brandy extract': 'apple flavoring',
    
    // Arabic alternatives
    'خنزير': 'لحم بقر حلال أو لحم غنم',
    'بيكون': 'بيكون بقر حلال أو ديك رومي',
    'هام': 'هام ديك رومي حلال',
    'كحول': 'عصير عنب أو عصير رمان',
    'خمر': 'عصير عنب',
    'نبيذ': 'عصير عنب',
    'جيلاتين': 'جيلاتين حلال أو أجار أجار'
};

// Function to check if an ingredient is halal
function isIngredientHalal(ingredient) {
    if (!ingredient || typeof ingredient !== 'string') {
        return true; // Default to halal if ingredient is not specified
    }
    
    const ingredientLower = ingredient.toLowerCase().trim();
    
    // Check against all non-halal categories
    for (const category in NON_HALAL_INGREDIENTS) {
        const nonHalalItems = NON_HALAL_INGREDIENTS[category];
        
        for (const item of nonHalalItems) {
            if (ingredientLower.includes(item.toLowerCase()) || 
                item.toLowerCase().includes(ingredientLower)) {
                return false;
            }
        }
    }
    
    return true;
}

// Function to get halal alternative for non-halal ingredient
function getHalalAlternative(ingredient) {
    if (!ingredient || typeof ingredient !== 'string') {
        return null;
    }
    
    const ingredientLower = ingredient.toLowerCase().trim();
    
    // Check for exact matches first
    if (HALAL_ALTERNATIVES[ingredientLower]) {
        return HALAL_ALTERNATIVES[ingredientLower];
    }
    
    // Check for partial matches
    for (const nonHalal in HALAL_ALTERNATIVES) {
        if (ingredientLower.includes(nonHalal.toLowerCase()) || 
            nonHalal.toLowerCase().includes(ingredientLower)) {
            return HALAL_ALTERNATIVES[nonHalal];
        }
    }
    
    return null;
}

// Function to filter halal recipes
function filterHalalRecipes(recipes) {
    if (!Array.isArray(recipes)) {
        return [];
    }
    
    return recipes.filter(recipe => {
        // Check if recipe is already marked as halal
        if (recipe.halal_certification === 'certified' || 
            recipe.halal_status === 'certified' ||
            recipe.is_halal === true) {
            return true;
        }
        
        // Check ingredients
        if (recipe.ingredients && Array.isArray(recipe.ingredients)) {
            for (const ingredient of recipe.ingredients) {
                const ingredientName = ingredient.item || ingredient.name || ingredient;
                if (!isIngredientHalal(ingredientName)) {
                    return false;
                }
            }
        }
        
        // Check recipe name and description for non-halal items
        const recipeName = recipe.meal_name || recipe.name || '';
        const recipeDescription = recipe.description || '';
        
        if (!isIngredientHalal(recipeName) || !isIngredientHalal(recipeDescription)) {
            return false;
        }
        
        return true;
    });
}

// Function to replace non-halal ingredients with halal alternatives
function replaceWithHalalAlternatives(recipe) {
    if (!recipe || typeof recipe !== 'object') {
        return recipe;
    }
    
    const modifiedRecipe = { ...recipe };
    const replacements = [];
    
    // Process ingredients
    if (modifiedRecipe.ingredients && Array.isArray(modifiedRecipe.ingredients)) {
        modifiedRecipe.ingredients = modifiedRecipe.ingredients.map(ingredient => {
            const ingredientName = ingredient.item || ingredient.name || ingredient;
            
            if (!isIngredientHalal(ingredientName)) {
                const alternative = getHalalAlternative(ingredientName);
                if (alternative) {
                    replacements.push({
                        original: ingredientName,
                        replacement: alternative
                    });
                    
                    // Update ingredient
                    if (ingredient.item) {
                        ingredient.item = alternative;
                    } else if (ingredient.name) {
                        ingredient.name = alternative;
                    } else {
                        return alternative;
                    }
                }
            }
            
            return ingredient;
        });
    }
    
    // Add replacement information to recipe
    if (replacements.length > 0) {
        modifiedRecipe.halal_replacements = replacements;
        modifiedRecipe.halal_modified = true;
    }
    
    return modifiedRecipe;
}

// Function to generate halal compliance report
function generateHalalComplianceReport(recipes) {
    if (!Array.isArray(recipes)) {
        return {
            total: 0,
            halal: 0,
            non_halal: 0,
            modified: 0,
            compliance_rate: 0
        };
    }
    
    let halalCount = 0;
    let nonHalalCount = 0;
    let modifiedCount = 0;
    
    recipes.forEach(recipe => {
        if (recipe.halal_certification === 'certified' || 
            recipe.halal_status === 'certified' ||
            recipe.is_halal === true) {
            halalCount++;
        } else {
            const halalRecipes = filterHalalRecipes([recipe]);
            if (halalRecipes.length > 0) {
                halalCount++;
            } else {
                const modifiedRecipe = replaceWithHalalAlternatives(recipe);
                if (modifiedRecipe.halal_modified) {
                    modifiedCount++;
                } else {
                    nonHalalCount++;
                }
            }
        }
    });
    
    const total = recipes.length;
    const complianceRate = total > 0 ? ((halalCount + modifiedCount) / total) * 100 : 0;
    
    return {
        total,
        halal: halalCount,
        non_halal: nonHalalCount,
        modified: modifiedCount,
        compliance_rate: Math.round(complianceRate * 100) / 100
    };
}

// Function to add halal filter to existing recipe filtering
function enhanceRecipeFilteringWithHalal(originalFilterFunction) {
    return function(recipes, clientData) {
        // First apply original filtering
        let filteredRecipes = originalFilterFunction(recipes, clientData);
        
        // Check if halal filtering is requested
        const restrictions = clientData.foodRestrictions || '';
        const isHalalRequested = restrictions.toLowerCase().includes('halal') || 
                               restrictions.toLowerCase().includes('حلال');
        
        if (isHalalRequested) {
            // Apply halal filtering
            filteredRecipes = filterHalalRecipes(filteredRecipes);
            
            // If no halal recipes found, try with alternatives
            if (filteredRecipes.length === 0) {
                const originalFiltered = originalFilterFunction(recipes, clientData);
                filteredRecipes = originalFiltered.map(recipe => 
                    replaceWithHalalAlternatives(recipe)
                ).filter(recipe => 
                    filterHalalRecipes([recipe]).length > 0
                );
            }
        }
        
        return filteredRecipes;
    };
}

// Function to display halal compliance information
function displayHalalComplianceInfo(recipes, containerId) {
    const report = generateHalalComplianceReport(recipes);
    const container = document.getElementById(containerId);
    
    if (!container) {
        console.warn('Halal compliance container not found:', containerId);
        return;
    }
    
    const complianceHTML = `
        <div class="halal-compliance-info alert alert-info mb-3">
            <h6><i class="fas fa-certificate me-2 text-success"></i>معلومات الامتثال الحلال</h6>
            <div class="row text-center">
                <div class="col-3">
                    <div class="compliance-stat">
                        <div class="stat-number text-success">${report.halal}</div>
                        <div class="stat-label">وصفات حلال</div>
                    </div>
                </div>
                <div class="col-3">
                    <div class="compliance-stat">
                        <div class="stat-number text-warning">${report.modified}</div>
                        <div class="stat-label">معدلة للحلال</div>
                    </div>
                </div>
                <div class="col-3">
                    <div class="compliance-stat">
                        <div class="stat-number text-danger">${report.non_halal}</div>
                        <div class="stat-label">غير حلال</div>
                    </div>
                </div>
                <div class="col-3">
                    <div class="compliance-stat">
                        <div class="stat-number text-primary">${report.compliance_rate}%</div>
                        <div class="stat-label">معدل الامتثال</div>
                    </div>
                </div>
            </div>
        </div>
    `;
    
    container.innerHTML = complianceHTML;
}

// Export functions for use in other modules
if (typeof module !== 'undefined' && module.exports) {
    module.exports = {
        isIngredientHalal,
        getHalalAlternative,
        filterHalalRecipes,
        replaceWithHalalAlternatives,
        generateHalalComplianceReport,
        enhanceRecipeFilteringWithHalal,
        displayHalalComplianceInfo,
        NON_HALAL_INGREDIENTS,
        HALAL_ALTERNATIVES
    };
}

// Global functions for browser environment
window.HalalFilter = {
    isIngredientHalal,
    getHalalAlternative,
    filterHalalRecipes,
    replaceWithHalalAlternatives,
    generateHalalComplianceReport,
    enhanceRecipeFilteringWithHalal,
    displayHalalComplianceInfo,
    NON_HALAL_INGREDIENTS,
    HALAL_ALTERNATIVES
};