# MrX's Inventory Tool

Mr. X owns a store that sells almost everything you think about. Now he wants a inventory management system to manage his inventory. Mr. X feels that controlling his inventory through SMS from his mobile will be revolutionary. So as a prequel, he decides that he wants a system that accepts one line commands and performs the respective operation.

Below is the list of commands he needs in the system:
* create itemName costPrice sellingPrice
       	Whenever Mr. X wants to add a new item to his store he issues a create command. This command creates a new item in the inventory with the given cost price and selling price. The prices are rounded off to two decimal places.

* delete itemName
      	If Mr. X decides not to sell an item anymore, then he simply issues a delete command. This command will remove the item from the inventory.

* updateBuy itemName quantity
      	Whenever Mr. X purchases additional quantity of the mentioned item, then he issues a updateBuy command. This command should increase the quantity of the mentioned item.

* updateSell itemName quantity
      	Whenever Mr. X sells some item, then he issues a updateSell command. This command should deduct the quantity of the mentioned item.

* report
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

# Building the tool 

You can build the binary with this command (OSX specific):

	$ docker-compose run -e GOOS=darwin --rm builder go build -o mrx-inventory

The magic test invocation

    $ docker run --rm -v "$PWD":/go/src/ -w /go/src golang:1.7 go test inventory -coverprofile coverage.out

The magic coverage invocation

    $ docker run --rm -v "$PWD":/go/src/ -w /go/src golang:1.7 go tool cover --html=coverage.out -o coverage.html 


