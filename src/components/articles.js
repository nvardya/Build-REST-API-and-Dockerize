import { React } from 'react'


    const Articles = ({ articles }) => {
      return (
        <div>
          <center><h1>Article Results</h1></center>
          {articles.map((article) => (
            <div class="card">
              <div class="card-body">
                <h5 class="card-title">{article.Title}</h5>
                <h6 class="card-subtitle mb-2 text-muted">{article.Abstract}</h6>
                <p class="card-text">{article.PublishedDate}</p>
                <p class="card-text">{article.URL}</p>
              </div>
            </div>
          ))}
        </div>
      );
    }

    export default Articles;
