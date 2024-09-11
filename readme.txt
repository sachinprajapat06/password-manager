Explanation:
Verifier: Verify the strength of your password.
Sighup: Creates a user and stores the hashed password.
Login: Authenticates a user by comparing the password with the hashed version.
Store Password: Encrypts the password using bcrypt and stores it in the database.
Get Password: Retrieves the stored encrypted password.

Note: Make sure to change your db information in controller file 