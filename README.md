Goat Robotics
------------------------
GoatRobotics is a real-time chat application developed using Go and the Gin Framework. The application allows clients to join a chat room, send messages, and view message history. It leverages concurrency tools such as sync.Map and RWMutex to handle multiple clients and messages efficiently, ensuring high performance even with numerous active users. The system exposes a RESTful API for communication, which includes features for message broadcasting, client management, and error handling.

Additionally, GoatRobotics provides a user-friendly static web UI that allows users to interact with the chat room seamlessly. You can access the UI at http://localhost:8080/static when the server is running. The project highlights the use of Go for building scalable, efficient, and real-time systems.

To Build Docker Image
--------------------------

docker build -t myapp .


To run Docker Container
--------------------------
    
docker run -d -p 8080:8080 --name myapp-container myapp

To Automate the both(in linux)
------------------------------------

./docker-script.sh
