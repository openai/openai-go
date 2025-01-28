// PMLL.go
// Persistent Memory with LRU Cache Implementation

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Logger Utility
type Logger struct {
	mutex sync.Mutex
}

// Log logs messages with different severity levels.
func (l *Logger) Log(level string, message string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	timestamp := time.Now().UTC().Format(time.RFC3339)
	switch strings.ToUpper(level) {
	case "INFO":
		fmt.Printf("[%s] [INFO] %s\n", timestamp, message)
	case "WARN":
		fmt.Printf("[%s] [WARN] %s\n", timestamp, message)
	case "ERROR":
		fmt.Printf("[%s] [ERROR] %s\n", timestamp, message)
	default:
		fmt.Printf("[%s] [UNKNOWN] %s\n", timestamp, message)
	}
}

// LRU Cache Implementation
type LRUCache struct {
	capacity int
	mutex    sync.Mutex
	cache    map[string]*listNode
	head     *listNode
	tail     *listNode
}

// listNode represents a node in the doubly linked list.
type listNode struct {
	key   string
	value string
	prev  *listNode
	next  *listNode
}

// NewLRUCache initializes a new LRU cache with the specified capacity.
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*listNode),
		head:     nil,
		tail:     nil,
	}
}

// Get retrieves a value from the cache.
// It returns the value and a boolean indicating if the key was found.
func (c *LRUCache) Get(key string) (string, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if node, exists := c.cache[key]; exists {
		c.moveToFront(node)
		return node.value, true
	}
	return "", false
}

// Put adds a key-value pair to the cache.
// If the key already exists, it updates the value and moves it to the front.
// If the cache exceeds its capacity, it evicts the least recently used item.
func (c *LRUCache) Put(key string, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if node, exists := c.cache[key]; exists {
		node.value = value
		c.moveToFront(node)
		return
	}

	newNode := &listNode{
		key:   key,
		value: value,
		prev:  nil,
		next:  nil,
	}

	if c.head == nil {
		c.head = newNode
		c.tail = newNode
	} else {
		newNode.next = c.head
		c.head.prev = newNode
		c.head = newNode
	}

	c.cache[key] = newNode

	if len(c.cache) > c.capacity {
		c.evict()
	}
}

// Clear purges all items from the cache.
func (c *LRUCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache = make(map[string]*listNode)
	c.head = nil
	c.tail = nil
}

// moveToFront moves a given node to the front of the list.
func (c *LRUCache) moveToFront(node *listNode) {
	if node == c.head {
		return
	}

	// Detach node
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}

	// Update tail if needed
	if node == c.tail {
		c.tail = node.prev
	}

	// Move to front
	node.prev = nil
	node.next = c.head
	if c.head != nil {
		c.head.prev = node
	}
	c.head = node

	// Update tail if list was only one node
	if c.tail == nil {
		c.tail = node
	}
}

// evict removes the least recently used item from the cache.
func (c *LRUCache) evict() {
	if c.tail == nil {
		return
	}
	delete(c.cache, c.tail.key)
	if c.tail.prev != nil {
		c.tail = c.tail.prev
		c.tail.next = nil
	} else {
		// Only one element was present
		c.head = nil
		c.tail = nil
	}
}

// PersistentMemory manages in-memory data with LRU caching and persistent storage.
type PersistentMemory struct {
	memoryFile     string
	cache          *LRUCache
	memoryData     map[string]string
	memoryVersions map[string][]string
	mutex          sync.Mutex
	logger         *Logger
}

// NewPersistentMemory initializes a new PersistentMemory instance.
func NewPersistentMemory(memoryFile string, cacheCapacity int, logger *Logger) *PersistentMemory {
	pm := &PersistentMemory{
		memoryFile:     memoryFile,
		cache:          NewLRUCache(cacheCapacity),
		memoryData:     make(map[string]string),
		memoryVersions: make(map[string][]string),
		logger:         logger,
	}

	pm.LoadMemory()
	return pm
}

// AddMemory adds or updates a memory entry asynchronously.
func (pm *PersistentMemory) AddMemory(key string, value string) {
	pm.mutex.Lock()
	pm.memoryData[key] = value
	pm.cache.Put(key, value)
	pm.logger.Log("INFO", fmt.Sprintf("Added/Updated memory entry for key: %s", key))
	pm.mutex.Unlock()

	go pm.SaveMemory()
}

// GetMemory retrieves a memory entry with caching.
func (pm *PersistentMemory) GetMemory(key string) (string, bool) {
	// Check cache first
	if value, found := pm.cache.Get(key); found {
		pm.logger.Log("INFO", fmt.Sprintf("Cache hit for key: %s", key))
		return value, true
	}

	// Lock to access memoryData
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	value, exists := pm.memoryData[key]
	if !exists {
		pm.logger.Log("WARN", fmt.Sprintf("Key not found: %s", key))
		return "", false
	}

	// Update cache
	pm.cache.Put(key, value)
	pm.logger.Log("INFO", fmt.Sprintf("Cache miss for key: %s. Loaded from storage.", key))

	return value, true
}

// ClearMemory clears all memory entries asynchronously.
func (pm *PersistentMemory) ClearMemory() {
	pm.mutex.Lock()
	pm.memoryData = make(map[string]string)
	pm.memoryVersions = make(map[string][]string)
	pm.cache.Clear()
	pm.logger.Log("INFO", "Cleared all memory entries and cache.")
	pm.mutex.Unlock()

	go pm.SaveMemory()
}

// AddMemoryVersion adds a new version to a memory entry asynchronously.
func (pm *PersistentMemory) AddMemoryVersion(key string, value string) {
	pm.mutex.Lock()
	pm.memoryVersions[key] = append(pm.memoryVersions[key], value)
	pm.memoryData[key] = value
	pm.cache.Put(key, value)
	version := len(pm.memoryVersions[key]) - 1
	pm.logger.Log("INFO", fmt.Sprintf("Added new memory version %d for key: %s", version, key))
	pm.mutex.Unlock()

	go pm.SaveMemory()
}

// GetMemoryVersion retrieves a specific version of a memory entry.
func (pm *PersistentMemory) GetMemoryVersion(key string, version int) (string, bool) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	versions, exists := pm.memoryVersions[key]
	if !exists {
		pm.logger.Log("WARN", fmt.Sprintf("No versions found for key: %s", key))
		return "", false
	}

	if version < 0 || version >= len(versions) {
		pm.logger.Log("WARN", fmt.Sprintf("Version %d out of range for key: %s", version, key))
		return "", false
	}

	value := versions[version]
	pm.logger.Log("INFO", fmt.Sprintf("Retrieved version %d for key: %s", version, key))
	return value, true
}

// LoadMemory loads memory data and versions from the JSON file.
func (pm *PersistentMemory) LoadMemory() {
	file, err := os.Open(pm.memoryFile)
	if err != nil {
		pm.logger.Log("WARN", "Memory file not found. Starting with empty memory.")
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	data := make(map[string]interface{})
	if err := decoder.Decode(&data); err != nil {
		pm.logger.Log("ERROR", fmt.Sprintf("Failed to parse memory file: %v", err))
		return
	}

	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Load memory_data
	if memData, ok := data["memory_data"].(map[string]interface{}); ok {
		for k, v := range memData {
			if strVal, ok := v.(string); ok {
				pm.memoryData[k] = strVal
				pm.cache.Put(k, strVal)
			}
		}
	}

	// Load memory_versions
	if memVersions, ok := data["memory_versions"].(map[string]interface{}); ok {
		for k, v := range memVersions {
			if versions, ok := v.([]interface{}); ok {
				for _, ver := range versions {
					if strVer, ok := ver.(string); ok {
						pm.memoryVersions[k] = append(pm.memoryVersions[k], strVer)
					}
				}
			}
		}
	}

	pm.logger.Log("INFO", "Memory loaded from file.")
}

// SaveMemory saves memory data and versions to the JSON file.
func (pm *PersistentMemory) SaveMemory() {
	pm.mutex.Lock()
	data := map[string]interface{}{
		"memory_data":     pm.memoryData,
		"memory_versions": pm.memoryVersions,
	}
	pm.mutex.Unlock()

	file, err := os.Create(pm.memoryFile)
	if err != nil {
		pm.logger.Log("ERROR", fmt.Sprintf("Failed to open memory file for writing: %v", err))
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		pm.logger.Log("ERROR", fmt.Sprintf("Failed to write memory file: %v", err))
		return
	}

	pm.logger.Log("INFO", "Memory saved to file.")
}

// CLI Interface
type CLI struct {
	pm        *PersistentMemory
	exitFlag  bool
	inputChan chan string
}

// NewCLI initializes a new CLI instance.
func NewCLI(pm *PersistentMemory) *CLI {
	return &CLI{
		pm:        pm,
		exitFlag:  false,
		inputChan: make(chan string),
	}
}

// Run starts the CLI loop.
func (cli *CLI) Run() {
	scanner := bufio.NewScanner(os.Stdin)

	for !cli.exitFlag {
		fmt.Println("\nAvailable Commands:")
		fmt.Println("1. Add")
		fmt.Println("2. Get")
		fmt.Println("3. Clear")
		fmt.Println("4. AddVersion")
		fmt.Println("5. GetVersion")
		fmt.Println("6. Exit")
		fmt.Print("Enter command number: ")

		if !scanner.Scan() {
			break
		}
		command := strings.TrimSpace(scanner.Text())

		switch command {
		case "1":
			cli.Add(scanner)
		case "2":
			cli.Get(scanner)
		case "3":
			cli.Clear()
		case "4":
			cli.AddVersion(scanner)
		case "5":
			cli.GetVersion(scanner)
		case "6":
			cli.Exit()
		default:
			fmt.Println("Invalid command. Please try again.")
		}
	}
}

// Add handles the Add command.
func (cli *CLI) Add(scanner *bufio.Scanner) {
	fmt.Print("Enter key: ")
	if !scanner.Scan() {
		return
	}
	key := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter value: ")
	if !scanner.Scan() {
		return
	}
	value := strings.TrimSpace(scanner.Text())

	cli.pm.AddMemory(key, value)
}

// Get handles the Get command.
func (cli *CLI) Get(scanner *bufio.Scanner) {
	fmt.Print("Enter key: ")
	if !scanner.Scan() {
		return
	}
	key := strings.TrimSpace(scanner.Text())

	if value, found := cli.pm.GetMemory(key); found {
		fmt.Printf("Value: %s\n", value)
	} else {
		fmt.Println("Key not found.")
	}
}

// Clear handles the Clear command.
func (cli *CLI) Clear() {
	cli.pm.ClearMemory()
	fmt.Println("All memory entries cleared.")
}

// AddVersion handles the AddVersion command.
func (cli *CLI) AddVersion(scanner *bufio.Scanner) {
	fmt.Print("Enter key: ")
	if !scanner.Scan() {
		return
	}
	key := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter version value: ")
	if !scanner.Scan() {
		return
	}
	value := strings.TrimSpace(scanner.Text())

	cli.pm.AddMemoryVersion(key, value)
}

// GetVersion handles the GetVersion command.
func (cli *CLI) GetVersion(scanner *bufio.Scanner) {
	fmt.Print("Enter key: ")
	if !scanner.Scan() {
		return
	}
	key := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter version number: ")
	if !scanner.Scan() {
		return
	}
	versionStr := strings.TrimSpace(scanner.Text())
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		fmt.Println("Invalid version number.")
		return
	}

	if value, found := cli.pm.GetMemoryVersion(key, version); found {
		fmt.Printf("Value at version %d: %s\n", version, value)
	} else {
		fmt.Println("Version not found.")
	}
}

// Exit handles the Exit command.
func (cli *CLI) Exit() {
	cli.exitFlag = true
	fmt.Println("Exiting CLI.")
}

// Main Function
func main() {
	// Initialize Logger
	logger := &Logger{}

	// Initialize PersistentMemory with a memory file and cache capacity
	pm := NewPersistentMemory("persistent_memory.json", 1000, logger)

	// Start CLI
	cli := NewCLI(pm)
	cli.Run()
}
// PMLL.go
// Persistent Memory with LRU Cache Implementation

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Logger Utility
// Provides thread-safe logging with different severity levels.
type Logger struct {
	mutex sync.Mutex
}

// Log logs messages with specified severity levels: INFO, WARN, ERROR.
func (l *Logger) Log(level string, message string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// Format the current timestamp in UTC.
	timestamp := time.Now().UTC().Format(time.RFC3339)

	// Log the message based on the severity level.
	switch strings.ToUpper(level) {
	case "INFO":
		fmt.Printf("[%s] [INFO] %s\n", timestamp, message)
	case "WARN":
		fmt.Printf("[%s] [WARN] %s\n", timestamp, message)
	case "ERROR":
		fmt.Printf("[%s] [ERROR] %s\n", timestamp, message)
	default:
		fmt.Printf("[%s] [UNKNOWN] %s\n", timestamp, message)
	}
}

// LRU Cache Implementation
// Manages in-memory caching with Least Recently Used eviction policy.
type LRUCache struct {
	capacity int                                     // Maximum number of items the cache can hold.
	mutex    sync.Mutex                              // Ensures thread-safe operations.
	cache    map[string]*listNode                    // Maps keys to their corresponding list nodes.
	head     *listNode                               // Most recently used item.
	tail     *listNode                               // Least recently used item.
}

// listNode represents a node in the doubly linked list used for the LRU cache.
type listNode struct {
	key   string        // The key of the cached item.
	value string        // The value of the cached item.
	prev  *listNode     // Pointer to the previous node in the list.
	next  *listNode     // Pointer to the next node in the list.
}

// NewLRUCache initializes a new LRU cache with the specified capacity.
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*listNode),
		head:     nil,
		tail:     nil,
	}
}

// Get retrieves a value from the cache.
// If the key exists, it moves the corresponding node to the front (most recently used).
// Returns the value and a boolean indicating whether the key was found.
func (c *LRUCache) Get(key string) (string, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check if the key exists in the cache map.
	if node, exists := c.cache[key]; exists {
		// Move the accessed node to the front to mark it as recently used.
		c.moveToFront(node)
		return node.value, true
	}
	return "", false
}

// Put adds a key-value pair to the cache.
// If the key already exists, it updates the value and moves the node to the front.
// If the cache exceeds its capacity, it evicts the least recently used item.
func (c *LRUCache) Put(key string, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check if the key already exists.
	if node, exists := c.cache[key]; exists {
		// Update the value and move the node to the front.
		node.value = value
		c.moveToFront(node)
		return
	}

	// Create a new node and add it to the front of the list.
	newNode := &listNode{
		key:   key,
		value: value,
		prev:  nil,
		next:  c.head,
	}

	if c.head != nil {
		c.head.prev = newNode
	}
	c.head = newNode

	// If the list was empty, set the tail to the new node.
	if c.tail == nil {
		c.tail = newNode
	}

	// Add the new node to the cache map.
	c.cache[key] = newNode

	// Evict the least recently used item if capacity is exceeded.
	if len(c.cache) > c.capacity {
		c.evict()
	}
}

// Clear purges all items from the cache.
func (c *LRUCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache = make(map[string]*listNode)
	c.head = nil
	c.tail = nil
}

// moveToFront moves a given node to the front of the doubly linked list.
func (c *LRUCache) moveToFront(node *listNode) {
	// If the node is already at the front, do nothing.
	if node == c.head {
		return
	}

	// Detach the node from its current position.
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}

	// If the node is the tail, update the tail pointer.
	if node == c.tail {
		c.tail = node.prev
	}

	// Move the node to the front.
	node.prev = nil
	node.next = c.head
	if c.head != nil {
		c.head.prev = node
	}
	c.head = node
}

// evict removes the least recently used item from the cache.
func (c *LRUCache) evict() {
	if c.tail == nil {
		return
	}
	// Remove the tail node from the cache map.
	delete(c.cache, c.tail.key)

	// Move the tail pointer to the previous node.
	if c.tail.prev != nil {
		c.tail = c.tail.prev
		c.tail.next = nil
	} else {
		// If there's only one node, reset head and tail.
		c.head = nil
		c.tail = nil
	}
}

// PersistentMemory manages in-memory data with LRU caching and persistent storage.
type PersistentMemory struct {
	memoryFile     string                              // Path to the JSON file for persistent storage.
	cache          *LRUCache                           // In-memory LRU cache.
	memoryData     map[string]string                   // Current memory entries.
	memoryVersions map[string][]string                 // Versioned memory entries.
	mutex          sync.Mutex                          // Ensures thread-safe access to memory data.
	logger         *Logger                             // Logger instance for logging operations.
}

// NewPersistentMemory initializes a new PersistentMemory instance.
// Parameters:
// - memoryFile: Path to the JSON file for storing memory data.
// - cacheCapacity: Maximum number of items the LRU cache can hold.
// - logger: Logger instance for logging.
func NewPersistentMemory(memoryFile string, cacheCapacity int, logger *Logger) *PersistentMemory {
	pm := &PersistentMemory{
		memoryFile:     memoryFile,
		cache:          NewLRUCache(cacheCapacity),
		memoryData:     make(map[string]string),
		memoryVersions: make(map[string][]string),
		logger:         logger,
	}

	pm.LoadMemory() // Load existing memory data from the file, if available.
	return pm
}

// AddMemory adds or updates a memory entry asynchronously.
// It updates the in-memory data and cache, then triggers an asynchronous save to persistent storage.
func (pm *PersistentMemory) AddMemory(key string, value string) {
	pm.mutex.Lock()
	pm.memoryData[key] = value
	pm.cache.Put(key, value)
	pm.logger.Log("INFO", fmt.Sprintf("Added/Updated memory entry for key: %s", key))
	pm.mutex.Unlock()

	// Asynchronously save memory data to the file.
	go pm.SaveMemory()
}

// GetMemory retrieves a memory entry with caching.
// It first checks the LRU cache; if not found, it fetches from persistent storage.
// Returns the value and a boolean indicating whether the key was found.
func (pm *PersistentMemory) GetMemory(key string) (string, bool) {
	// Check cache first.
	if value, found := pm.cache.Get(key); found {
		pm.logger.Log("INFO", fmt.Sprintf("Cache hit for key: %s", key))
		return value, true
	}

	// Lock to access memoryData.
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	value, exists := pm.memoryData[key]
	if !exists {
		pm.logger.Log("WARN", fmt.Sprintf("Key not found: %s", key))
		return "", false
	}

	// Update cache with the retrieved value.
	pm.cache.Put(key, value)
	pm.logger.Log("INFO", fmt.Sprintf("Cache miss for key: %s. Loaded from storage.", key))

	return value, true
}

// ClearMemory clears all memory entries and the cache asynchronously.
// It also triggers an asynchronous save to persist the cleared state.
func (pm *PersistentMemory) ClearMemory() {
	pm.mutex.Lock()
	pm.memoryData = make(map[string]string)
	pm.memoryVersions = make(map[string][]string)
	pm.cache.Clear()
	pm.logger.Log("INFO", "Cleared all memory entries and cache.")
	pm.mutex.Unlock()

	// Asynchronously save the cleared state to the file.
	go pm.SaveMemory()
}

// AddMemoryVersion adds a new version to a memory entry asynchronously.
// It appends the new version to the memoryVersions map and updates the current memoryData.
// Then, it triggers an asynchronous save to persistent storage.
func (pm *PersistentMemory) AddMemoryVersion(key string, value string) {
	pm.mutex.Lock()
	pm.memoryVersions[key] = append(pm.memoryVersions[key], value)
	pm.memoryData[key] = value
	pm.cache.Put(key, value)
	version := len(pm.memoryVersions[key]) - 1
	pm.logger.Log("INFO", fmt.Sprintf("Added new memory version %d for key: %s", version, key))
	pm.mutex.Unlock()

	// Asynchronously save memory data to the file.
	go pm.SaveMemory()
}

// GetMemoryVersion retrieves a specific version of a memory entry.
// Returns the value and a boolean indicating whether the version was found.
func (pm *PersistentMemory) GetMemoryVersion(key string, version int) (string, bool) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	versions, exists := pm.memoryVersions[key]
	if !exists {
		pm.logger.Log("WARN", fmt.Sprintf("No versions found for key: %s", key))
		return "", false
	}

	if version < 0 || version >= len(versions) {
		pm.logger.Log("WARN", fmt.Sprintf("Version %d out of range for key: %s", version, key))
		return "", false
	}

	value := versions[version]
	pm.logger.Log("INFO", fmt.Sprintf("Retrieved version %d for key: %s", version, key))
	return value, true
}

// LoadMemory loads memory data and versions from the JSON file.
// If the file does not exist, it starts with empty memory.
func (pm *PersistentMemory) LoadMemory() {
	file, err := os.Open(pm.memoryFile)
	if err != nil {
		pm.logger.Log("WARN", "Memory file not found. Starting with empty memory.")
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	data := make(map[string]interface{})
	if err := decoder.Decode(&data); err != nil {
		pm.logger.Log("ERROR", fmt.Sprintf("Failed to parse memory file: %v", err))
		return
	}

	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Load memory_data from the JSON.
	if memData, ok := data["memory_data"].(map[string]interface{}); ok {
		for k, v := range memData {
			if strVal, ok := v.(string); ok {
				pm.memoryData[k] = strVal
				pm.cache.Put(k, strVal)
			}
		}
	}

	// Load memory_versions from the JSON.
	if memVersions, ok := data["memory_versions"].(map[string]interface{}); ok {
		for k, v := range memVersions {
			if versions, ok := v.([]interface{}); ok {
				for _, ver := range versions {
					if strVer, ok := ver.(string); ok {
						pm.memoryVersions[k] = append(pm.memoryVersions[k], strVer)
					}
				}
			}
		}
	}

	pm.logger.Log("INFO", "Memory loaded from file.")
}

// SaveMemory saves memory data and versions to the JSON file.
// It locks the memory data during the save operation to ensure consistency.
func (pm *PersistentMemory) SaveMemory() {
	pm.mutex.Lock()
	data := map[string]interface{}{
		"memory_data":     pm.memoryData,
		"memory_versions": pm.memoryVersions,
	}
	pm.mutex.Unlock()

	file, err := os.Create(pm.memoryFile)
	if err != nil {
		pm.logger.Log("ERROR", fmt.Sprintf("Failed to open memory file for writing: %v", err))
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // For pretty-printing.
	if err := encoder.Encode(data); err != nil {
		pm.logger.Log("ERROR", fmt.Sprintf("Failed to write memory file: %v", err))
		return
	}

	pm.logger.Log("INFO", "Memory saved to file.")
}

// CLI Interface
// Provides a command-line interface for interacting with the PersistentMemory system.
type CLI struct {
	pm        *PersistentMemory // Reference to the PersistentMemory instance.
	exitFlag  bool               // Indicates whether to exit the CLI loop.
	inputChan chan string        // Channel for handling asynchronous inputs (optional).
}

// NewCLI initializes a new CLI instance.
func NewCLI(pm *PersistentMemory) *CLI {
	return &CLI{
		pm:        pm,
		exitFlag:  false,
		inputChan: make(chan string),
	}
}

// Run starts the CLI loop, presenting available commands and handling user input.
func (cli *CLI) Run() {
	scanner := bufio.NewScanner(os.Stdin)

	for !cli.exitFlag {
		// Display available commands.
		fmt.Println("\nAvailable Commands:")
		fmt.Println("1. Add")
		fmt.Println("2. Get")
		fmt.Println("3. Clear")
		fmt.Println("4. AddVersion")
		fmt.Println("5. GetVersion")
		fmt.Println("6. Exit")
		fmt.Print("Enter command number: ")

		// Read user input.
		if !scanner.Scan() {
			break
		}
		command := strings.TrimSpace(scanner.Text())

		// Handle commands based on user input.
		switch command {
		case "1":
			cli.Add(scanner)
		case "2":
			cli.Get(scanner)
		case "3":
			cli.Clear()
		case "4":
			cli.AddVersion(scanner)
		case "5":
			cli.GetVersion(scanner)
		case "6":
			cli.Exit()
		default:
			fmt.Println("Invalid command. Please try again.")
		}
	}
}

// Add handles the "Add" command.
// Prompts the user for a key and value, then adds the memory entry.
func (cli *CLI) Add(scanner *bufio.Scanner) {
	fmt.Print("Enter key: ")
	if !scanner.Scan() {
		return
	}
	key := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter value: ")
	if !scanner.Scan() {
		return
	}
	value := strings.TrimSpace(scanner.Text())

	cli.pm.AddMemory(key, value)
}

// Get handles the "Get" command.
// Prompts the user for a key and retrieves the corresponding memory entry.
func (cli *CLI) Get(scanner *bufio.Scanner) {
	fmt.Print("Enter key: ")
	if !scanner.Scan() {
		return
	}
	key := strings.TrimSpace(scanner.Text())

	if value, found := cli.pm.GetMemory(key); found {
		fmt.Printf("Value: %s\n", value)
	} else {
		fmt.Println("Key not found.")
	}
}

// Clear handles the "Clear" command.
// Clears all memory entries and the cache.
func (cli *CLI) Clear() {
	cli.pm.ClearMemory()
	fmt.Println("All memory entries cleared.")
}

// AddVersion handles the "AddVersion" command.
// Prompts the user for a key and a versioned value, then adds the version.
func (cli *CLI) AddVersion(scanner *bufio.Scanner) {
	fmt.Print("Enter key: ")
	if !scanner.Scan() {
		return
	}
	key := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter version value: ")
	if !scanner.Scan() {
		return
	}
	value := strings.TrimSpace(scanner.Text())

	cli.pm.AddMemoryVersion(key, value)
}

// GetVersion handles the "GetVersion" command.
// Prompts the user for a key and version number, then retrieves the specific version.
func (cli *CLI) GetVersion(scanner *bufio.Scanner) {
	fmt.Print("Enter key: ")
	if !scanner.Scan() {
		return
	}
	key := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter version number: ")
	if !scanner.Scan() {
		return
	}
	versionStr := strings.TrimSpace(scanner.Text())
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		fmt.Println("Invalid version number.")
		return
	}

	if value, found := cli.pm.GetMemoryVersion(key, version); found {
		fmt.Printf("Value at version %d: %s\n", version, value)
	} else {
		fmt.Println("Version not found.")
	}
}

// Exit handles the "Exit" command.
// Sets the exit flag to true, causing the CLI loop to terminate.
func (cli *CLI) Exit() {
	cli.exitFlag = true
	fmt.Println("Exiting CLI.")
}

// Main Function
// Initializes the PersistentMemory system and starts the CLI.
func main() {
	// Initialize Logger.
	logger := &Logger{}

	// Initialize PersistentMemory with:
	// - "persistent_memory.json" as the storage file.
	// - 1000 as the cache capacity.
	pm := NewPersistentMemory("persistent_memory.json", 1000, logger)

	// Initialize and run the CLI.
	cli := NewCLI(pm)
	cli.Run()
}
