# Coffee Machine Simulator

### Features
- Supports multiple Mixers(Outlets) to prepare dishes in parallel. This is achieved by running multiple go routines that takes requests from a single queue and generates prepared item to a different queue
- All mixers use a common device stock inventory that is thread safe using locks
- Machine supports warning indicator that takes a list of ingrident quantity threshold and matches it with current stock
- A Mixer(Outlet) can only produce a item as long as it has enough ingridents to do so in device stock

### Instructions
- Create a stock store
    - InmemoryStock is availabe in code
    - Store can be initialized with ingridents at start
- Create a device that takes
    - Warning thresholds
    - Stock store
- FillIngridents() in device store if needed
- CreateOutlet(Mixer) as needed for this device
    - Takes a list of items that this mixer can cook
    - User requests are taken from a channel
    - Generates prepared dishes to a output channel
- Use RequestItem() for requesting for dishes
- CheckRefill() can be used monitor if devices needs to be filled again

### How to run

- Checkout repo
- Should have go and maketool installed
- Run tests
```
make test
```
- Run sample binary
```
make && ./dunzo
```