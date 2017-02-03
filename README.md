# MrX's Inventory Tool

Mr. X owns a store that sells almost everything you think about. Now he wants a inventory management system to manage his inventory. Mr. X feels that controlling his inventory through SMS from his mobile will be revolutionary. So as a prequel, he decides that he wants a system that accepts one line commands and performs the respective operation.

Below is the list of commands he needs in the system:
* `create itemName costPrice sellingPrice`
       	Whenever Mr. X wants to add a new item to his store he issues a create command. This command creates a new item in the inventory with the given cost price and selling price. The prices are rounded off to two decimal places.

* `delete itemName`
      	If Mr. X decides not to sell an item anymore, then he simply issues a delete command. This command will remove the item from the inventory.

* `updateBuy itemName quantity`
      	Whenever Mr. X purchases additional quantity of the mentioned item, then he issues a updateBuy command. This command should increase the quantity of the mentioned item.

* `updateSell itemName quantity`
      	Whenever Mr. X sells some item, then he issues a updateSell command. This command should deduct the quantity of the mentioned item.

* `report`
      	Whenever Mr. X wants to view his inventory list he issues the report command. This command should print the current inventory details in the specified format sorted by alphabetical order. Apart from printing the inventory it has to report on the profit made by Mr. X since last report generation. Where profit is calculated by:  `sum(sellingPrice - costPrice)` of the sold items multiplied by `(no. of items sold - costPrice)` of the deleted items.


Sample Input
```
create Book01 10.50 13.79
create Food01 1.47 3.98
create Med01 30.63 34.29
create Tab01 57.00 84.98
updateBuy Tab01 100
updateSell Tab01 2
updateBuy Food01 500
updateBuy Book01 100
updateBuy Med01 100
updateSell Food01 1
updateSell Food01 1
updateSell Tab01 2
report
delete Book01
updateSell Tab01 5
create Mobile01 10.51 44.56
updateBuy Mobile01 250
updateSell Food01 5
updateSell Mobile01 4
updateSell Med01 10
report
#
```
Expected Output
```
              	INVENTORY REPORT
Item Name 	Bought At    	Sold At       	AvailableQty    	Value
--------- 	---------    	-------       	-----------     	-------
Book01    	10.50          	13.79               	100    	1050.00
Food01     	1.47           	3.98               	498     	732.06
Med01     	30.63          	34.29               	100    	3063.00
Tab01     	57.00          	84.98                	96    	5472.00
---------------------------------------------------------------------------
Total value                                                     	10317.06
Profit since previous report                                      	116.94


              	INVENTORY REPORT
Item Name 	Bought At    	Sold At  	AvailableQty    	Value
--------- 	---------    	-------  	-----------     	-------
Food01          	1.47      	3.98       	493           	724.71
Med01          	30.63     	34.29        	90          	2756.70
Mobile01       	10.51     	44.56       	246          	2585.46
Tab01          	57.00     	84.98        	91          	5187.00
---------------------------------------------------------------------------
Total value                                                   	11253.87
Profit since previous report                                   	-724.75
```

# Development
The only requirement for building this application is docker, preferably on some flavor of 'nix.

## Building the tool

You can build the binary with this command on OSX:

	$ docker run --rm -v "$PWD":/go/src/ -w /go/src -e GOOS=darwin golang:1.7 go build -o mrx-inventory

If you are building the tool on standard linux, use this command instead:

	$ docker run --rm -v "$PWD":/go/src/ -w /go/src golang:1.7 go build -o mrx-inventory

## Running the tests
This is the command to run all of the unit tests, and get the coverage statistics

	$ docker run --rm -v "$PWD":/go/src/ -w /go/src golang:1.7 go test inventory -coverprofile coverage.out

This is the command to convert the coverage statistic into a friendly HTML report

	$ docker run --rm -v "$PWD":/go/src/ -w /go/src golang:1.7 go tool cover --html=coverage.out -o coverage.html 

You can now open `coverage.html` in your browser

# Usage
The mrx-inventory program is a daemon, that will listen for commands on TCP port 8333 by default.  You can start it like this

	$ ./mrx-inventory

This will constantly display the last report.  In order to interact with the daemon, a tool like netcat will be very useful  For example, you could interact with the daemon like this:

	$ echo -n "create BaconNinjaThing 15.00 204.38" | nc localhost 8333
	$ echo -n "create MiniPizza 0.03 1.01" | nc localhost 8333
	$ echo -n "report" | nc localhost 8333

The daemon will then present a report like this:

```
Name                    Quantity  Bought At    Sold At    Bought Value
-----                      -----      -----      -----           -----
BaconNinjaThing                0      15.00     204.38            0.00
MiniPizza                      0       0.03       1.01            0.00
-----                      -----      -----      -----           -----
                              Inventory Value                     0.00
                              Revenue Since Last Report           0.00
                              Cost Since Last Report              0.00
                              Profit Since Last Report            0.00
```

# Design
The problem statement opens with a vision for having a system controlled via SMS.  In order to hit this goal, there are a few issues that need to be hammered out right away.  A working data model needs to be created, so that is possible to decide how commands should interact.  Second, the command parser needs to be developed, so that it is possible to interact with the system.  Third, the goal of using SMS presents several concurrency issues, with the possibility of commands being lost or arriving out of order.  The data model needs to handle these concerns directly.

## Data Model
The problem states four types of commands that can change the state of the inventory, and one command that changes the state of accumulated cash.  In order to test the effects these commands have on the system, the state was modeled explicitly.  The most important interface was the StateEntry interface, which you can see the key pieces below:

	type StateEntry interface {
		NextState(accum State) (State, error)
		//Human Readability Methods Omitted
	}

The interface would be consumed in the following manner:

	state, err = entry.NextState(state)

Since the state is passed in and returned explicitly, it allows the modifications the specific entry makes to be tested down to the smallest detail.  This turned out to be very useful, as I had made a mistake in my original interpretation of the spec.  You can see git commit `7d68125ddd86d3` in particular to watch as I fix this specific mistake.

Also, by having errors being explicitly returned, it is easier to test situations that are strictly valid.  For example, what happens when an entry is created twice?  What happens when the updateBuy or updateSell order is for a negative quantity?  What happens when you try to sell more inventory that you currently have?  Each of these can easily be unit tested by the virtue of making state explicit.

## Parsing
Parsing the data is easy to state, difficult to get right.  The majority of the parsing is done in a function called ParseLine, it has the following signature

    ParseLine(line string,/*otherDeps*/) (StateEntry,error)

You would invoke the parser like so

    entry, err = ParseLine(line,/*otherDeps)

When it is possible to parse the line, a concrete type the implements the StateEntry interface is returned.  However, which implementation of the interface that is returned is hidden from the user. You'll notice that there also is an error type returned as well.  If the function was not able to successfully parse the line, an error value will be returned. 

The parser is deliberately very strict. There was a choice between making the parse stricter but more predictable, or friendlier, but perhaps less predictable.  The decision was made to error on the side of predictability, because a predictable system is required for (eventually) delivering a friendlier system.

## Concurrency
Finally, the system will eventually be required to accept input from SMS.  This means there are several concurrency issues that need to be evaluated.  The system will need to be able to accept input from multiple sources concurrently, and the messages will arrive potentially out of order.

It was the concurrently arriving messages issues that drove me to deliver this prototype using go (well, and it was fun).  Go provides several features to make this easy, specifically first class internal processes, and channel support for communicating between these processes.

The application has one internal process the is in charge of managing state.  It receives inputs on an internal channel, which it then uses to update the internal application state.  Interacting with this internal process is simple, and all of the complexity is hidden from clients of this process by the channel.  All the consuming process needs to do is send a message on the channel, and a response will be returned to original process once the state management process is complete.

If you look in the application entry point, `main.go`, you'll see that the message quickly gets from the TCP perimeter to the appropriate channel.

Other major problem, messages arriving out of order, is a little trickier.  The application doesn't call of any transaction or batching features that would normally be used to mitigate risks like this.  The only mitigate at this point is to thoroughly test the data model.  For example, you'll see code to protect `updateBuy` or `updateSell` commands from being executed if the items doesn't exist yet.  There are many other safeguards like this one.

One last concurrency concern that was address was with the `report` command.  The was a design concern about the report command taking a long time, and providing an opportunity for a race condition while commands are still being processed(I/O can be slow).  In order to preserve the system throughput, a copy of the appropriate state object is sent to a separate async goroutine, in order to keep everything thread safe.  If this becomes a performance bottleneck, an immutable version of the state object could be investigated.
