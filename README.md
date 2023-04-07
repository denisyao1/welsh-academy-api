# Welsh Academy API
Welsh Academy API is a demo API dedicated to provide recipes to cheddar lovers around the world.

## installation
To install welsh API on your local computer, follow these steps:

1. Install git on your computer

1. Install docker and docker compose on your computer. Generally docker compose is include in docker desktop. So you can install docker desktop on your machine and get both

1. Clone the project by running this command `https://github.com/denisyao1/welsh-academy-api.git` or `git@github.com:denisyao1/welsh-academy-api.git`

1. Navigate to the root of the folder and create a .env file according to the env.sample file.
1. Run the command `docker-compose up`

1. Open your browser and enter url : http://localhost:3000/docs ; you should see the swagger documentation of the API and voilà welsh-academy-api is running on your machine

## Usage
The swagger documentation of the API shows you how to use it.

The base URL of the API is http://localhost:3000/api/v1 .
A **default admin user** is created when the API is built. It’s **username is admin so as its password**.
**It's recommended to change the default admin password**. 
You can do it on the swagger page by log in using the default admin credential and then change its password by doing a PACTH request on /users/password-change directly on swagger page. The swagger page describes all http request you can perform with the API and the inputs and / or parameters each request can accept.
Many endpoints need authentication to be accessible.
**Welsh API save token in http cookies so you don't need to fill manually token in request header to use it**.

You can also create new users by providing their username, password and specifying if has admin privilege or not.
A user can know its username and role (isAdmin) by making a GET request on /users/my-infos.

A user can :
- list all existing ingredients 
- list all possible recipes (with or without ingredient constraints); to do so he must add ingredients name's  as request parameter
- flag/unflag recipes as his favorite ones
- list his favorite recipes

A admin user can do all thing a normal user can do plus :
- Create a new admin or normal user
- Create ingredients : to create an ingredient it must provide only its name.
- Create recipes of meals using the previously created ingredients : to create a recipe, he must provide the recipe **name**, the recipe **making** and the list of the **name of ingredients** of recipe.

Contact us if you have any suggestion or question.
You hope you will enjoy the API.
