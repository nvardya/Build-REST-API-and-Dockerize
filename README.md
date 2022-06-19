# Buiild a Web App and Deploy with Docker

This project is a web application that is deployed with a microservice archutecture. Docker Containers and Docker Compose support this architecture

![Docker Diagram](https://user-images.githubusercontent.com/53916435/174461271-08c9ab87-9140-4833-9f38-c1e2a1af454e.jpg)

This web appplication utilized the following applications for build:
  1. A Node.js application to scrape data from The New York Time's REST API endpoint and send the data to a PostgreSQL database
  2. A Golang application to build a REST API (GET, PUT, POST operations) to interact with the PostgreSQL database
  3. A React application to prompt users to search from the PostgreSQL database via an API request on the Golang application

![REST API Diagram1](https://user-images.githubusercontent.com/53916435/174461411-1c310bdb-421c-48ce-9484-1bf5b8e2e099.jpg)

#Table of Contents
1. Overview
2. Build Node.js Web Scraper
  2a. ScrapeFromNYT.js
3. Build Golang REST API
  3a. Install gorilla/mux
4. Build React Frontend
5. Build Dockefiles and Docker Compose file


# 1. Build Node.js Web Scraper
The Node.js part of the web application will perform the following steps:
1) Access the New York Time's API endpoint and make a GET call to gather the most popular articles from the past 7 days
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

# 2. Build Golang REST API
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
# 3. Build React Frontend
React will be used to interact with the user. There are 2 React components that will be created for this piece of the application:
1) `App.js`: This will prompt the user to enter and search for a keyword that they are looking for. If this key word exists in one of the article records in our PostgreSQL table, it will be displayed
2) `article.js`: This will return the article results from the PostgreSQL table based off the keyword input entered by the user

The most important part of building `App.js` is to include the **useState** React Hook. This allows us to return the appropriate articles from the PostgreSQL datababse whenever the user enters a new keyword in the seearch bar

```javascript
const [getResult, setGetResult] = useState(null);
```
# 3. Build DockerFiles and Docker Compose YAML file
As you know at this point, this application is really broken down into 3 separate microservices:

1) Node.JS Web Scraper
2) Golang REST API
3) React Frontend

Because they are all somehwat independent of each other, I decided to create Docker images and subsequently containers for each of them. Within each directory of this repository, you will find a `Dockerfile` file. This file is what will be used to create the Docker image. However, I wanted to take it a step further. Instead of having to make 3 separate commands to launch Docker containers for each part of my application, I decided to create a `docker-compose.yml` file that will orchestrate and luanch all 3 of my containers. See the below code (configs) in `docker-compose.yml`:

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
Remmeber, a Dockerfile will need to be created for each directory/microservice in this web application. You can see that all 3 Dockerfiles are mentioned in the above YAML code. Run `docker-compose.yml` in your terminal to launch the YAML file and your application will be launched!

# Getting Started with Create React App

This project was bootstrapped with [Create React App](https://github.com/facebook/create-react-app).


### `npm start`

Runs the app in the development mode.\
Open [http://localhost:3000](http://localhost:3000) to view it in your browser.

The page will reload when you make changes.\
You may also see any lint errors in the console.

### `npm test`

Launches the test runner in the interactive watch mode.\
See the section about [running tests](https://facebook.github.io/create-react-app/docs/running-tests) for more information.

### `npm run build`

Builds the app for production to the `build` folder.\
It correctly bundles React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.\
Your app is ready to be deployed!

See the section about [deployment](https://facebook.github.io/create-react-app/docs/deployment) for more information.

### `npm run eject`

**Note: this is a one-way operation. Once you `eject`, you can't go back!**

If you aren't satisfied with the build tool and configuration choices, you can `eject` at any time. This command will remove the single build dependency from your project.

Instead, it will copy all the configuration files and the transitive dependencies (webpack, Babel, ESLint, etc) right into your project so you have full control over them. All of the commands except `eject` will still work, but they will point to the copied scripts so you can tweak them. At this point you're on your own.

You don't have to ever use `eject`. The curated feature set is suitable for small and middle deployments, and you shouldn't feel obligated to use this feature. However we understand that this tool wouldn't be useful if you couldn't customize it when you are ready for it.

## Learn More

You can learn more in the [Create React App documentation](https://facebook.github.io/create-react-app/docs/getting-started).

To learn React, check out the [React documentation](https://reactjs.org/).

### Code Splitting

This section has moved here: [https://facebook.github.io/create-react-app/docs/code-splitting](https://facebook.github.io/create-react-app/docs/code-splitting)

### Analyzing the Bundle Size

This section has moved here: [https://facebook.github.io/create-react-app/docs/analyzing-the-bundle-size](https://facebook.github.io/create-react-app/docs/analyzing-the-bundle-size)

### Making a Progressive Web App

This section has moved here: [https://facebook.github.io/create-react-app/docs/making-a-progressive-web-app](https://facebook.github.io/create-react-app/docs/making-a-progressive-web-app)

### Advanced Configuration

This section has moved here: [https://facebook.github.io/create-react-app/docs/advanced-configuration](https://facebook.github.io/create-react-app/docs/advanced-configuration)

### Deployment

This section has moved here: [https://facebook.github.io/create-react-app/docs/deployment](https://facebook.github.io/create-react-app/docs/deployment)

### `npm run build` fails to minify

This section has moved here: [https://facebook.github.io/create-react-app/docs/troubleshooting#npm-run-build-fails-to-minify](https://facebook.github.io/create-react-app/docs/troubleshooting#npm-run-build-fails-to-minify)
