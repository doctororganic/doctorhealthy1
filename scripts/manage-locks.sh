#!/bin/bash

# Lock file management script for multi-developer coordination
# Usage: ./scripts/manage-locks.sh [command] [args]

LOCKS_DIR=".locks"
LOCK_EXTENSION=".lock"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Create locks directory if it doesn't exist
mkdir -p "$LOCKS_DIR"

# Function to display help
show_help() {
    echo "Lock File Management System"
    echo ""
    echo "Usage: $0 [command] [args]"
    echo ""
    echo "Commands:"
    echo "  lock <file> <developer_id> <description> <eta>  Create a lock file"
    echo "  unlock <file>                                      Remove a lock file"
    echo "  check                                              Check all active locks"
    echo "  check <file>                                       Check if specific file is locked"
    echo "  clean                                              Remove stale locks (older than 24h)"
    echo "  help                                               Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 lock main.go DEV3 \"Adding new API endpoint\" \"30min\""
    echo "  $0 unlock main.go"
    echo "  $0 check"
    echo "  $0 check main.go"
    echo "  $0 clean"
}

# Function to create a lock file
create_lock() {
    local file="$1"
    local dev_id="$2"
    local description="$3"
    local eta="$4"
    
    if [ -z "$file" ] || [ -z "$dev_id" ] || [ -z "$description" ] || [ -z "$eta" ]; then
        echo -e "${RED}Error: Missing arguments${NC}"
        echo "Usage: $0 lock <file> <developer_id> <description> <eta>"
        return 1
    fi
    
    local lock_file="$LOCKS_DIR/${file}${LOCK_EXTENSION}"
    
    if [ -f "$lock_file" ]; then
        echo -e "${RED}Error: File $file is already locked${NC}"
        echo -e "${YELLOW}Current lock:$(cat "$lock_file")${NC}"
        return 1
    fi
    
    echo "$dev_id: $description - ETA: $eta" > "$lock_file"
    echo -e "${GREEN}✅ Lock created for $file${NC}"
    echo -e "${BLUE}Lock details: $dev_id: $description - ETA: $eta${NC}"
}

# Function to remove a lock file
remove_lock() {
    local file="$1"
    
    if [ -z "$file" ]; then
        echo -e "${RED}Error: Missing file argument${NC}"
        echo "Usage: $0 unlock <file>"
        return 1
    fi
    
    local lock_file="$LOCKS_DIR/${file}${LOCK_EXTENSION}"
    
    if [ ! -f "$lock_file" ]; then
        echo -e "${RED}Error: No lock found for $file${NC}"
        return 1
    fi
    
    echo -e "${YELLOW}Removing lock for $file${NC}"
    echo -e "${BLUE}Previous lock: $(cat "$lock_file")${NC}"
    rm "$lock_file"
    echo -e "${GREEN}✅ Lock removed for $file${NC}"
}

# Function to check all locks
check_all_locks() {
    echo -e "${BLUE}Checking all active locks...${NC}"
    echo ""
    
    local found_locks=false
    
    for lock_file in "$LOCKS_DIR"/*"$LOCK_EXTENSION"; do
        if [ -f "$lock_file" ]; then
            local file_name=$(basename "$lock_file" "$LOCK_EXTENSION")
            local lock_content=$(cat "$lock_file")
            local lock_time=$(stat -c %Y "$lock_file" 2>/dev/null || stat -f %m "$lock_file" 2>/dev/null)
            local current_time=$(date +%s)
            local age=$((current_time - lock_time))
            local age_hours=$((age / 3600))
            local age_minutes=$(((age % 3600) / 60))
            
            echo -e "${YELLOW}🔒 Locked: $file_name${NC}"
            echo -e "   ${BLUE}Details: $lock_content${NC}"
            echo -e "   ${BLUE}Age: ${age_hours}h ${age_minutes}m${NC}"
            echo ""
            
            found_locks=true
        fi
    done
    
    if [ "$found_locks" = false ]; then
        echo -e "${GREEN}✅ No active locks found${NC}"
    fi
}

# Function to check specific file lock
check_file_lock() {
    local file="$1"
    
    if [ -z "$file" ]; then
        echo -e "${RED}Error: Missing file argument${NC}"
        echo "Usage: $0 check <file>"
        return 1
    fi
    
    local lock_file="$LOCKS_DIR/${file}${LOCK_EXTENSION}"
    
    if [ -f "$lock_file" ]; then
        local lock_content=$(cat "$lock_file")
        local lock_time=$(stat -c %Y "$lock_file" 2>/dev/null || stat -f %m "$lock_file" 2>/dev/null)
        local current_time=$(date +%s)
        local age=$((current_time - lock_time))
        local age_hours=$((age / 3600))
        local age_minutes=$(((age % 3600) / 60))
        
        echo -e "${YELLOW}🔒 File $file is locked${NC}"
        echo -e "   ${BLUE}Details: $lock_content${NC}"
        echo -e "   ${BLUE}Age: ${age_hours}h ${age_minutes}m${NC}"
        return 0
    else
        echo -e "${GREEN}✅ File $file is not locked${NC}"
        return 1
    fi
}

# Function to clean stale locks
clean_stale_locks() {
    echo -e "${BLUE}Checking for stale locks (older than 24 hours)...${NC}"
    
    local current_time=$(date +%s)
    local stale_age=$((24 * 3600)) # 24 hours in seconds
    local found_stale=false
    
    for lock_file in "$LOCKS_DIR"/*"$LOCK_EXTENSION"; do
        if [ -f "$lock_file" ]; then
            local lock_time=$(stat -c %Y "$lock_file" 2>/dev/null || stat -f %m "$lock_file" 2>/dev/null)
            local age=$((current_time - lock_time))
            
            if [ $age -gt $stale_age ]; then
                local file_name=$(basename "$lock_file" "$LOCK_EXTENSION")
                local lock_content=$(cat "$lock_file")
                local age_hours=$((age / 3600))
                local age_minutes=$(((age % 3600) / 60))
                
                echo -e "${YELLOW}🗑️  Removing stale lock: $file_name${NC}"
                echo -e "   ${BLUE}Details: $lock_content${NC}"
                echo -e "   ${BLUE}Age: ${age_hours}h ${age_minutes}m${NC}"
                
                rm "$lock_file"
                found_stale=true
            fi
        fi
    done
    
    if [ "$found_stale" = false ]; then
        echo -e "${GREEN}✅ No stale locks found${NC}"
    else
        echo -e "${GREEN}✅ Stale locks cleaned${NC}"
    fi
}

# Main command router
case "$1" in
    "lock")
        create_lock "$2" "$3" "$4" "$5"
        ;;
    "unlock")
        remove_lock "$2"
        ;;
    "check")
        if [ -n "$2" ]; then
            check_file_lock "$2"
        else
            check_all_locks
        fi
        ;;
    "clean")
        clean_stale_locks
        ;;
    "help"|"--help"|"-h"|"")
        show_help
        ;;
    *)
        echo -e "${RED}Error: Unknown command '$1'${NC}"
        echo ""
        show_help
        exit 1
        ;;
esac