# SEO Web Scraper
#### Video Demo:  TBD
#### Description:
SEO Web scraper built in Go.
Quickly look up your websites and check their SEO settings.
###### The app checks the following:
- 'robots' META tag
- 'keywords' META tag
- 'description' META tag
- if the website has its sitemap linked onto the page
- Open Graph METAs
- Twitter Card METAs
- Dublin Core METAs
- JSON-LD schema
- Facebook Pixel integration
- all external links found onto the page
###### Teck stack:
- Go - https://go.dev/
- Echo web framework - https://github.com/labstack/echo
- Colly scraper - https://github.com/gocolly/colly
- SQLite - https://github.com/mattn/go-sqlite3
- JWT - https://pkg.go.dev/github.com/golang-jwt/jwt/v5
- Bcrypt - https://pkg.go.dev/golang.org/x/crypto/bcrypt
- HTMX - https://github.com/bigskysoftware/htmx
- Go templ - https://github.com/a-h/templ/
- godotenv - https://github.com/joho/godotenv

### App presentation:
#### Front-end
Using Go Templ I've created a layout that applies to all pages (/views/layout) which renders the HTML structure and imports the HTMX script and the global CSS file. The navigation has items for homepage, login and register (if the user is not logged in) and history (if the user is logged in).
- Homepage allows you to scrape a website. If you are logged in it will automatically save the searched website into the database so you can see past searches when accessing the History page.
- Login allows you to login as an user
- Register allows you to register as an user
- History renders the websites you've searched in the past
The search, login and register forms send data to different API endpoints with the help of HTMX.
#### Back-end
In the app's entry point (main.go) the 'godotenv' library is loaded which contains variables for the server's listen address, the database name and the JWT secret. Then a SQLite storage is created, this does the following: opens the database, creates tables for users and websites, seeds the database with dummy users. Lastly, the API server is created and started.
##### SQLite storage
Using receiver functions (methods in the OOP world), the following database operations are created:
- Create user
- Get user by email
- Update user
- Delete user 
- Create website
- Get websites by user ID
##### API Server
The API has the following endpoints:
- GET "/" - Renders home page
- GET "/login" - Renders login page
- GET "/register" - Renders register page
- GET "/history" - Renders history page, user must be logged in to access this page
- POST "/register" - Grabs the form values sent, encrypts the password using bcrypt, stores the user into the database, renders HTML partials for operation successful or failed.
- POST "/login" - Grabs the form values sent, gets the user from database based on the email, compares the hashed passwords using bcrypt, generates a JWT token, sets a cookie with the generated token string and redirects to homepage
- POST "/website" - Grabs the website URL sent through the form, creates a new scraper based on the website's URL and checks the SEO settings. If the user is logged in, the JWT token string is grabbed from the cookie and validated, the website is stored into the database and the HTML with all scraped info is rendered onto the page.
### Next steps:
- Switch to MySQL database
- Add functionality to scrape the whole website, not just one page
- Add functionality to check for broken (404) links 
- Add functionality to send email with the full report
- Add functionality to schedule website for scraping daily and send notification to user
