#Fampay Assignment

##Architecture Diagram
![](https://i.ibb.co/55QqLQH/Architecture-02.png)


###Local Setup

- make docker-build-api
- make docker-build-worker

The default config assumes that mongodb is running at localhost:27017 with proper indexing enabled.
Please follow mongo_setup.txt present in the repository to setup collections and indexes.

The APIKey configured in the default.toml will automatically be written to mongodb on api startup. 


#### For Kafka Local Setup:
- Follow the steps for your operating system. For Ubuntu 20.04 you can follow [this link](https://tecadmin.net/how-to-install-apache-kafka-on-ubuntu-20-04/)

The default config again assumes that kafka is running locally at localhost:9092. Please make changes to the config if the ports are different.


#### For local testing:
- Please use the provided postman collection present in the repo to test locally.


#### Future Enhancements:
- Make the Worker shutdown gracefully.
- provide auto configuration for mutiple kafka workers. Please note multiple workers can still be run now with config changes.
- Add unit and integ test cases. The API has been tested locally but testcases are not written currently.
- Add github CI support for builds.
