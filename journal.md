
## Week 1 (26 mar - 31 mar)

### First Meeting (26 March 2018)

To start with the we decided to go with the Starbucks online application. Together we discussed the team tasks that we will be doing in coming weeks. For the first week would be working on our application and design layout. We are planning to meet again on saturday i. e. 31st march to decide the next steps.

### March 30, 2018
* We completed basic frontend UI template pages for the starbucks application. we each are going to work and complete individual pages for different feartures like order page, update order page, locations, etc.
* We are also discussing what architecture and deployment strategies to use in order to achieve high AKF cube scaling

### April 2, 2018
We read through the article: https://www.benefitfocus.com/blogs/design-engineering/architecture-cubed and began to decide the architecture which we are going to implement. Apart from that we decided to use NodeJS as part of our frontend which will be deployed on EC2 instances along with load balancer and auto scailing group.


### April 7, 2018
We worked on the basic design html pages. We created the starbucks menu page and added the other features of qty and size while selecting the drink. Also implemented the logic for calculating the price and setting the store location.

### April 12, 2018
* We discussed about the backend requirement and what REST backend we should use.
* We have decided to build a go based REST API to use in this project.


### April 14, 2018
* Added few Go REST APIs and checked transaction flow between frontend and backend.
* Need to also add a Redis DB to check the whole flow from front to back till we have completed out personal projects to connect the cluster.

### April 16, 2018
* Modified more Go REST api to handle more functionality

### April 18, 2018
* tried connection GO with Dynomite on redis with 5 nodes
* some testing still needed.

### April 19, 2018
* consulted with our final design with Pranav, the TA, and he has approved our design.

### April 21, 2018
* Presented our deployment strategy to Professor for achieving high AKF cube scaling but he has suggested few changes which we are considering making.

### April 22, 2018
* Have finalized which API we need to support to achive maximum efficiency and support AKF scaling.
* url will be in the format of {store_location}/order to enable sharding.
* the following API will be supported
* Place Order: 		POST	- 	starbucks/order
* Get order: 		GET 	- 	starbucks/order/order_id
* Get orders:		GET 	- 	starbucks/orders
* Delete order: 	DEL 	- 	starbucks/order/order_id
* Update order: 	PUT 	- 	starbucks/order/order_id
* Pay for order:	POST 	- 	starbucks/order/order_id/pay
* making changes to the go backend accordingly.

### April 23, 2018
* Discussed with TA to get a clear picture of Y-axis scaling and how to apply it to our structure

### April 25, 2018
* Deployed the frontend Node.js files on heroku
* need to deploy the GO app in AWS to talk to the frontend. need a bit of testing

### April 27, 2018
* deployed Go application with load balancer on my vpc
* Frontend was able to talk to backend

### April 28, 2018
* Discussed further with the professor about our deployment architecture and he suggested few modifications to achieve higher level of AKF scaling
* We need to seperate our micro services and have them running together on different VPCs.

### May 02 , 2018
Discussed the pending tasks and the follow through actions to be taken due to some changes in the basic design strategy. Decided to meet on 3rd and 4th to complete the remaining tasks.

## May 03, 2018
Worked on individual go apis and completed the front end.
