# Project Name: Startbucks Online Application

## CMPE 281 Team - Phoenix


## Team Members:
* Animesh Grover
* Neha Yawalkar
* Watcharit Maharutainont


### Technologies Used:
* AWS Web Services
* Heroku (Frond-end deployment)
* Docker
* Kong
* Cassandra
* GoLang
* Nodejs
* Redis

### About Project :

* We have designed this application on the basis of AKF scaling. 
* The front end web server is deployed on the heroku. 
* Kong and cassandra are deployed on a separate VPC. The sharding will happen here based on the location of the store (SF, SJ,MV).
* Each team member has implemented a microservice in golang (Product/Cart/ Order) and these services are running on all of our VPCs on seperate instances.
* The instances are load balanced and attached to an auto-scaling group.
* 5 nodes of redis are running on each machine. The nodes are also configured to handle network partition.


### Architecture Diagram

![Architecture](https://github.com/nguyensjsu/team281-phoenix/blob/master/documents/AKF.png)

