# 🎓 BEGINNER'S DEPLOYMENT GUIDE
## Step-by-Step Instructions for Complete Beginners

**Time Required:** 30 minutes  
**Difficulty:** Easy (just follow the steps!)

---

## 📍 WHERE ARE YOU NOW?

You're on your Mac, and you have a folder called `nutrition-platform` somewhere (probably on Desktop).

---

## 🎯 PART 1: OPEN TERMINAL

### Step 1: Find Terminal
1. Press `Command (⌘) + Space` on your keyboard
2. Type: `terminal`
3. Press `Enter`
4. A black/white window will open - this is Terminal!

### Step 2: Navigate to Your Project
In Terminal, type these commands (press Enter after each):

```bash
# Go to Desktop (if your project is there)
cd Desktop

# Go into the project folder
cd nutrition-platform

# Check you're in the right place (should show files)
ls
```

**What you should see:** A list of files including `backend`, `frontend-nextjs`, `docker-compose.yml`, etc.

---

## 🎯 PART 2: RUN THE DEPLOYMENT SCRIPT

### Step 3: Make Script Executable
Copy and paste this into Terminal, then press Enter:

```bash
chmod +x DEPLOY-TO-COOLIFY-NOW.sh
```

**What this does:** Makes the script runnable (like double-clicking an app)

### Step 4: Run the Script
Copy and paste this into Terminal, then press Enter:

```bash
./DEPLOY-TO-COOLIFY-NOW.sh
```

**What happens:**
- Script generates secure passwords
- Creates deployment files
- Shows you credentials to save

### Step 5: Save the Credentials
The script will show something like:

```
📋 SAVE THESE CREDENTIALS SECURELY:
====================================
DB_PASSWORD=07352a3b890942733b2106ca142be5a0...
REDIS_PASSWORD=7913e6d120029f00714361f92888ddab...
JWT_SECRET=51a4aad6a61eb4c5e8f410b4517a6269...
```

**IMPORTANT:** 
1. Select all this text with your mouse
2. Press `Command (⌘) + C` to copy
3. Open Notes app (`Command + Space`, type "notes")
4. Press `Command (⌘) + V` to paste
5. Save the note as "Nutrition Platform Credentials"

---

## 🎯 PART 3: OPEN COOLIFY DASHBOARD

### Step 6: Open Coolify in Browser
1. Open Safari or Chrome
2. Go to: `https://api.doctorhealthy1.com`
3. You'll see the Coolify login page

### Step 7: Login to Coolify
- Enter your Coolify username
- Enter your Coolify password
- Click "Login"

**You should see:** Coolify dashboard with menu on the left

---

## 🎯 PART 4: CREATE PROJECT IN COOLIFY

### Step 8: Create New Project

**Click these in order:**

1. **Left Menu** → Click "Projects"
2. **Top Right** → Click "+ New Project" button
3. **Fill in the form:**
   - Name: `nutrition-platform`
   - Description: `AI-powered nutrition and health management platform`
4. **Click** "Create"

**You should see:** Your new project in the list

---

## 🎯 PART 5: ADD YOUR APPLICATION

### Step 9: Create New Resource

1. **Click** on your project name (`nutrition-platform`)
2. **Click** "+ New Resource" button
3. **Choose** "Docker Compose"

### Step 10: Configure Docker Compose

**You'll see a form. Fill it in:**

#### Basic Settings:
- **Name:** `nutrition-app`
- **Description:** `Main application`

#### Source:
- **Choose:** "Docker Compose"
- **Compose File:** Click "Upload" or "Paste"

### Step 11: Upload Docker Compose File

**Option A: Upload File**
1. Click "Upload File"
2. Navigate to: `Desktop/nutrition-platform/`
3. Select: `docker-compose.coolify.yml`
4. Click "Open"

**Option B: Paste Content**
1. Click "Paste Content"
2. Go back to Terminal
3. Type: `cat docker-compose.coolify.yml`
4. Copy all the text that appears
5. Paste into Coolify

---

## 🎯 PART 6: ADD ENVIRONMENT VARIABLES

### Step 12: Set Environment Variables

1. **In Coolify**, scroll down to "Environment Variables"
2. **Click** "+ Add Variable" button

**Now add these variables one by one:**

Go back to Terminal and type:
```bash
cat .env.coolify.secure
```

**You'll see something like:**
```
DB_HOST=postgres
DB_PORT=5432
DB_NAME=nutrition_platform
...
```

**For EACH line:**
1. Click "+ Add Variable"
2. **Name:** Copy the part BEFORE the `=` (e.g., `DB_HOST`)
3. **Value:** Copy the part AFTER the `=` (e.g., `postgres`)
4. Click "Add"

**Repeat for ALL variables** (about 20 variables)

**Important variables to add:**
- `DB_HOST` = `postgres`
- `DB_PORT` = `5432`
- `DB_NAME` = `nutrition_platform`
- `DB_USER` = `nutrition_user`
- `DB_PASSWORD` = (the long password from credentials)
- `REDIS_HOST` = `redis`
- `REDIS_PORT` = `6379`
- `REDIS_PASSWORD` = (the long password from credentials)
- `JWT_SECRET` = (the long secret from credentials)
- `PORT` = `8080`
- `ENVIRONMENT` = `production`
- `DOMAIN` = `super.doctorhealthy1.com`

---

## 🎯 PART 7: CONFIGURE DOMAINS

### Step 13: Add Domains

1. **Scroll down** to "Domains" section
2. **Click** "+ Add Domain"

**Add these domains:**

**Domain 1 (Frontend):**
- Domain: `super.doctorhealthy1.com`
- Port: `3000`
- Click "Add"

**Domain 2 (Backend):**
- Domain: `api.super.doctorhealthy1.com`
- Port: `8080`
- Click "Add"

---

## 🎯 PART 8: DEPLOY!

### Step 14: Start Deployment

1. **Scroll to top** of the page
2. **Click** the big "Deploy" button
3. **Wait** (this takes 5-10 minutes)

**What you'll see:**
- Build logs scrolling
- Progress indicators
- Status changing from "Building" → "Running"

### Step 15: Monitor Deployment

**Watch the logs:**
- Green text = Good ✅
- Red text = Error ❌
- Yellow text = Warning ⚠️

**Common messages (all normal):**
- "Pulling image..."
- "Building..."
- "Starting container..."
- "Health check passed"

---

## 🎯 PART 9: VERIFY IT WORKS

### Step 16: Check Health

**In Coolify:**
1. Look for "Status" indicator
2. Should show: 🟢 Running

**In Browser:**
1. Open new tab
2. Go to: `https://api.super.doctorhealthy1.com/health`
3. Should see: `{"status":"ok"}`

### Step 17: Access Your App

**Open these URLs:**

1. **Frontend:** `https://super.doctorhealthy1.com`
   - Should see: Your nutrition platform homepage

2. **Backend API:** `https://api.super.doctorhealthy1.com/api/v1/info`
   - Should see: JSON with API information

---

## 🎯 TROUBLESHOOTING

### Problem: "Script not found"
**Solution:**
```bash
# Make sure you're in the right folder
pwd
# Should show: /Users/yourname/Desktop/nutrition-platform

# If not, navigate there:
cd ~/Desktop/nutrition-platform
```

### Problem: "Permission denied"
**Solution:**
```bash
chmod +x DEPLOY-TO-COOLIFY-NOW.sh
```

### Problem: "Docker Compose file not found"
**Solution:**
```bash
# Check if file exists
ls -la docker-compose.coolify.yml

# If not, the script creates it when you run it
./DEPLOY-TO-COOLIFY-NOW.sh
```

### Problem: Coolify shows "Build Failed"
**Solution:**
1. Click on the failed deployment
2. Read the error message
3. Usually it's a missing environment variable
4. Go back and add the missing variable

### Problem: "Cannot connect to database"
**Solution:**
1. Check environment variables in Coolify
2. Make sure `DB_PASSWORD` matches the one from credentials
3. Restart the deployment

---

## 📱 QUICK REFERENCE

### Terminal Commands You Used:
```bash
cd Desktop/nutrition-platform          # Navigate to project
chmod +x DEPLOY-TO-COOLIFY-NOW.sh     # Make script runnable
./DEPLOY-TO-COOLIFY-NOW.sh            # Run deployment script
cat .env.coolify.secure                # View credentials
cat docker-compose.coolify.yml         # View Docker config
```

### Coolify URLs:
- **Dashboard:** https://api.doctorhealthy1.com
- **Your App:** https://super.doctorhealthy1.com
- **API:** https://api.super.doctorhealthy1.com

### Important Files:
- `DEPLOY-TO-COOLIFY-NOW.sh` - Deployment script
- `.env.coolify.secure` - Your credentials
- `docker-compose.coolify.yml` - Docker configuration
- `coolify.json` - Coolify settings

---

## 🎓 WHAT EACH STEP DID

1. **Terminal:** Command-line interface to run scripts
2. **chmod +x:** Made script executable
3. **./script.sh:** Ran the script
4. **Credentials:** Secure passwords for database/services
5. **Coolify Project:** Container for your app
6. **Docker Compose:** Instructions for running all services
7. **Environment Variables:** Configuration settings
8. **Domains:** URLs where your app is accessible
9. **Deploy:** Build and start your application

---

## ✅ SUCCESS CHECKLIST

- [ ] Terminal opened
- [ ] Navigated to project folder
- [ ] Ran deployment script
- [ ] Saved credentials in Notes
- [ ] Logged into Coolify
- [ ] Created project
- [ ] Added Docker Compose
- [ ] Added all environment variables
- [ ] Added domains
- [ ] Clicked Deploy
- [ ] Deployment shows "Running"
- [ ] Health check returns OK
- [ ] Can access frontend
- [ ] Can access backend API

---

## 🆘 NEED HELP?

### Check Logs in Coolify:
1. Click on your application
2. Click "Logs" tab
3. Look for error messages

### Check Status:
1. In Coolify dashboard
2. Look for 🟢 (running) or 🔴 (stopped)

### Restart Application:
1. Click on application
2. Click "Restart" button
3. Wait for it to start

---

## 🎉 YOU'RE DONE!

Your nutrition platform is now live at:
- **https://super.doctorhealthy1.com**

You can:
- Track nutrition
- Plan meals
- Manage workouts
- Monitor progress

**Congratulations! 🎊**

---

**Need to deploy again?** Just run:
```bash
cd ~/Desktop/nutrition-platform
./DEPLOY-TO-COOLIFY-NOW.sh
```

**Questions?** Check the logs in Coolify or review this guide!
