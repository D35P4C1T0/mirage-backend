Sure, here’s a basic structure of possible API methods and endpoints for your REST backend:

### Authentication

1. **Register**
    - Endpoint: `/api/auth/register`
    - Method: `POST`
    - Description: Register a new user account.

2. **Login**
    - Endpoint: `/api/auth/login`
    - Method: `POST`
    - Description: Authenticate a user and return a token.

3. **Logout**
    - Endpoint: `/api/auth/logout`
    - Method: `POST`
    - Description: Invalidate the user token.

### User Management

4. **Get User Profile**
    - Endpoint: `/api/users/{userId}`
    - Method: `GET`
    - Description: Retrieve user profile information.

5. **Update User Profile**
    - Endpoint: `/api/users/{userId}`
    - Method: `PUT`
    - Description: Update user profile information.

6. **Delete User**
    - Endpoint: `/api/users/{userId}`
    - Method: `DELETE`
    - Description: Remove a user account.

### Album Management

7. **Upload Album**
    - Endpoint: `/api/albums`
    - Method: `POST`
    - Description: Upload a new album (pack of pictures).

8. **Get Albums**
    - Endpoint: `/api/albums`
    - Method: `GET`
    - Description: Retrieve a list of albums.

9. **Get Album**
    - Endpoint: `/api/albums/{albumId}`
    - Method: `GET`
    - Description: Retrieve specific album details.

10. **Update Album**
    - Endpoint: `/api/albums/{albumId}`
    - Method: `PUT`
    - Description: Update an album's information.

11. **Delete Album**
    - Endpoint: `/api/albums/{albumId}`
    - Method: `DELETE`
    - Description: Remove an album.

### Picture Management

12. **Upload Picture**
    - Endpoint: `/api/albums/{albumId}/pictures`
    - Method: `POST`
    - Description: Upload a new picture to an album.

13. **Get Pictures**
    - Endpoint: `/api/albums/{albumId}/pictures`
    - Method: `GET`
    - Description: Retrieve a list of pictures in an album.

14. **Get Picture**
    - Endpoint: `/api/albums/{albumId}/pictures/{pictureId}`
    - Method: `GET`
    - Description: Retrieve specific picture details.

15. **Delete Picture**
    - Endpoint: `/api/albums/{albumId}/pictures/{pictureId}`
    - Method: `DELETE`
    - Description: Remove a picture from an album.

### Smart Frame Integration

16. **Send Album to Smart Frame**
    - Endpoint: `/api/smart-frames/{frameId}/albums`
    - Method: `POST`
    - Description: Send an album to a specified smart frame.

### AI Person Recognition (Future Implementation)

17. **Run Person Recognition**
    - Endpoint: `/api/albums/{albumId}/recognize`
    - Method: `POST`
    - Description: Run AI-based person recognition on an album's pictures.

18. **Get Recognition Results**
    - Endpoint: `/api/albums/{albumId}/recognition-results`
    - Method: `GET`
    - Description: Retrieve person recognition results for an album.
