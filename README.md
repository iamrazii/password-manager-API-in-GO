🔐 Passman API – Password Manager Backend:
 
 Password manager API built with Go (gin framework) and Microsoft SQL Server. It supports user authentication via JWT tokens and enables users to securely store, retrieve, update, and delete their passwords — all on protected routes.





 
 ⚙️ Features

. User registration & login

. JWT-based authentication

. Middleware for token validation & user context

. Store, retrieve, update, and delete passwords per user

. Clean and RESTful API structure

. Secure DB queries using sql.Named parameters

. A Dokcer File 






🔑 JWT Authentication Flow



1.User logs in and receives a JWT token.


2.The frontend must attach the token in the Authorization header for protected routes:
   
    Authorization: Bearer <token>


3.Middleware verifies this token, extracts the user_id, and injects it into the context for route handlers.


4.All operations like password creation or retrieval are scoped to the authenticated user.
