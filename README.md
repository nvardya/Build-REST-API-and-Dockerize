# Build a Microservice Based Web App and Deploy with Docker

This project is a web application built with a microservices archutecture. Docker Containers and Docker Compose support the deployment of these microservices. This document will outline the steps needed to build and deploy this web application.

# 1. Overview
3 total microservices were created to support the web application. A Docker Container was created for each of them. Docker Compose orchestrates the deploymenmt of all 3 containers:

![image](https://user-images.githubusercontent.com/53916435/174485385-0805fe92-ccc4-4732-b119-b61f10385b53.png)

This web appplication utilized the following for build:

  1. A Node.js application to scrape data from The New York Time's REST API endpoint and send the data to a PostgreSQL database
  2. A Golang application to build a REST API (GET, PUT, POST operations) to interact with the PostgreSQL database
  3. A React application to prompt users to search from the PostgreSQL database via an API request on the Golang application

![image](https://user-images.githubusercontent.com/53916435/174485404-7e723d09-c1b0-4722-9c71-ee791e596373.png)

# 2. Build Node.js Web Scraper
The Node.js Web Scraper will perform the following steps:

1) Access The New York Time's API endpoint and make an API call to gather the most popular articles from the past 7 days
2) Insert the articles gathered from the previous step into a PostgreSQL database

A `ScrapeFromNYT.js` file will be created.

A New York Times developer account will be required. Details of NYT's API endpoints can be found [here](https://developer.nytimes.com). Once you create your NYT API account, setting up the connection to your NYT API account in `ScrapeFromNYT.js` is quite simple:

```javascript
const URL_Setup = {
  hostname: 'api.nytimes.com',
  port: 443,
  path: '/svc/mostpopular/v2/emailed/1.json?api-key=********************',
  method: 'GET',
};
});
```

You will also need to setup a connection to your PostgreSQL database:

```javascript
//A connection pool was used as opposed to a client
const pool = new Pool({
  user: 'username',
  host: 'nyt-d********.us-east-2.rds.amazonaws.com',
  database: '****',
  password: '*****',
  port: 5432,
});
```

Please see `ScrapeFromNYT.js` for details on inserting data returned from NYT into the PostgreSQL table

# 3. Build Golang REST API
Golang will be used to create a REST API and connect it to the PostgreSQL database. GET, PUT, and POST operations will be created.

A `main.go` file will be created.

The `gorilla/mux` package is the foundation of creating the REST API. `gorilla/mux` implements a request router and matches incoming HTTP requests agaianst a list of registered routes. Adiditional details on the package can be found [here](https://github.com/gorilla/mux).

The followng code in `main.go` summarizes what will be built in this file (a request router and GET/PUT/POST handlers):

```golang
router := mux.NewRouter().StrictSlash(true)
router.HandleFunc("/article/{query}", GETHandler).Methods("GET")
router.HandleFunc("/article/{id}", PUTHandler).Methods("PUT")
router.HandleFunc("/article/insert", POSTHandler)
http.Handle("/", router)
log.Fatal(http.ListenAndServe(":10000", nil))
```
# 4. Build React Frontend
React will be used to interact with the user. There are 2 React components that will be created for this piece of the application:
1) `App.js`: This will prompt the user to enter and search for a keyword that they are looking for. If this key word exists in one of the article records in our PostgreSQL table, it will be displayed
2) `article.js`: This will return the article results from the PostgreSQL table based off the keyword input entered by the user

The most important part of building `App.js` is to include the **useState** React Hook. This allows us to return the appropriate articles from the PostgreSQL datababse whenever the user enters a new keyword in the seearch bar

```javascript
const [getResult, setGetResult] = useState(null);
```
# 5. Build DockerFiles and Docker Compose YAML file
As you know at this point, this application is really broken down into 3 separate microservices:

1) Node.JS Web Scraper
2) Golang REST API
3) React Frontend

Because they are all somehwat independent of each other, I decided to create Docker images and subsequently containers for each of them. Within each directory of this repository, you will find a `Dockerfile` file. This file is what will be used to create a Docker image of the microservice. However, I wanted to take it a step further. Instead of having to make 3 separate commands to launch Docker containers for each part of my application, I decided to create a `docker-compose.yml` file that will orchestrate and luanch all 3 of my containers. See the below code (configs) in `docker-compose.yml`:

```YAML
version: '3'
services:
  nytscrape:
    build: ./nytScrape
    volumes:
     - .:/code
  buildrestapi:
    build: ./buildrestapi
    ports:
      - "10000:10000"
  react-api:
    build: ./react-api
    ports:
      - "3000:3000"
```
Remmeber, a Dockerfile will need to be created for each directory/microservice in this web application. You can see that all 3 Dockerfiles are mentioned in the above YAML code. Additionally, I recommend downloading the Docker desktop app as you can see your Docker Images and Containers. Below is a snapshot from my Docker desktop app of my Docker Compose container:

<img width="1268" alt="Screen Shot 2022-06-18 at 9 17 13 PM" src="https://user-images.githubusercontent.com/53916435/174465644-3c0f137b-751a-4a9c-ac70-0bda50e43711.png">

Run `docker-compose up` in your terminal to launch the YAML file and your application will be launched!

https://user-images.githubusercontent.com/53916435/174465985-36ca9342-6346-46f8-aba3-b83573802a43.mov
