const https = require('https');
const { Pool } = require("pg");

const URL_Setup = {
  hostname: 'api.nytimes.com',
  port: 443,
  path: '/svc/mostpopular/v2/emailed/1.json?api-key=********************',
  method: 'GET',
};

//A connection pool was used as opposed to a client
const pool = new Pool({
  user: 'username',
  host: 'nyt-d********.us-east-2.rds.amazonaws.com',
  database: '****',
  password: '*****',
  port: 5432,
});

async function poolConnection() {
  const pool = new Pool(credentials);
  const now = await pool.query("SELECT NOW()");
  await pool.end();

  return now;
}

async function insertArticle(article) {
  const text = `
    INSERT INTO AllArticles (Title, Abstract, PublishedDate, URL)
    VALUES ($1, $2, $3, $4)
  `;
  const values = [article.Title, article.Abstract, article.PublishedDate, article.URL];
  return pool.query(text, values);
}

https
.get("https://api.nytimes.com/svc/mostpopular/v2/emailed/1.json?api-key=KoGc9hZD19iarsljrC6YsrsCQyYMtVXO", resp => {
    let data = "";

    //Create a string from the NYT API response body
    resp.on("data", returnString => {
      data += returnString;
    });

    // Parse the string into a JSON structure. Neeed this in order to loop and
    // pull from the respone
   resp.on("end", () => {
     let json_response = JSON.parse(data);
     let article_list =  json_response.results;
     for (var x in article_list) {

  // Insert article detail details into Postgres Table
        const newArticle = insertArticle({
          Title: article_list[x].title,
          Abstract: article_list[x].abstract,
          PublishedDate: article_list[x].published_date,
          URL: article_list[x].url
        });
     }
   });
})
.on("error", err => {
    console.log("Error: " + err.message);
  });
