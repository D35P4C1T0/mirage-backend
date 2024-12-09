package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// User Represents a user in the system
type User struct {
	ID               primitive.ObjectID   `bson:"_id,omitempty"`
	Username         string               `bson:"username" binding:"required"`    // Required unique username
	Email            string               `bson:"email" binding:"required,email"` // Required unique email
	Password         string               `bson:"password" binding:"required"`    // Required password
	UserProfileID    primitive.ObjectID   `bson:"user_profile_id,omitempty"`      // Links to the user's profile
	ProfilePictureID primitive.ObjectID   `bson:"profile_picture_id,omitempty"`   // Links to the user's profile picture
	AlbumsID         []primitive.ObjectID `bson:"albums_id,omitempty"`            // List of owned albums
	CreatedAt        time.Time            `bson:"created_at"`                     // User account creation timestamp
	UpdatedAt        time.Time            `bson:"updated_at"`                     // Last profile update timestamp
}

// UserProfile Represents a user's detailed profile
type UserProfile struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `bson:"name" binding:"required"`    // First name
	Surname string             `bson:"surname" binding:"required"` // Last name
	Sex     string             `bson:"sex"`                        // Gender
	DOB     time.Time          `bson:"dob"`                        // Date of birth
}

// ProfilePicture Represents a profile picture
type ProfilePicture struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ImageData []byte             `bson:"image_data,omitempty"`       // Binary profile picture data
	UserID    primitive.ObjectID `bson:"user_id" binding:"required"` // User associated with this picture
	CreatedAt time.Time          `bson:"created_at"`                 // When the profile picture was uploaded
}

// Album Represents an album owned by a user
type Album struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty"`
	Title         string               `bson:"title" binding:"required"`  // Required album title
	Description   string               `bson:"description,omitempty"`     // Optional description
	OwnerID       primitive.ObjectID   `bson:"user_id"`                   // Album owner
	TargetUserIDs []primitive.ObjectID `bson:"target_user_ids,omitempty"` // List of users it is shared with
	Tags          []string             `bson:"tags,omitempty"`            // Tags for categorization
	IsPrivate     bool                 `bson:"is_private"`                // Privacy setting
	CreatedAt     time.Time            `bson:"created_at"`                // Creation timestamp
	UpdatedAt     time.Time            `bson:"updated_at"`                // Last updated timestamp
}

// RecognizedFace Represents face recognition metadata
type RecognizedFace struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`  // Unique ID
	PictureID    primitive.ObjectID `bson:"picture_id"`     // Associated picture
	Name         string             `bson:"name,omitempty"` // Name (if available)
	Sex          string             `bson:"sex,omitempty"`  // Gender
	Confidence   float64            `bson:"confidence"`     // Recognition confidence level
	RecognizedAt time.Time          `bson:"recognized_at"`  // Time of recognition
	Age          int                `bson:"age,omitempty"`  // Estimated age
}

// Picture Represents a picture in an album
type Picture struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty"`
	PictureDataID primitive.ObjectID   `bson:"picture_data_id"`       // Picture data reference
	Thumbnail     []byte               `bson:"thumbnail,omitempty"`   // Thumbnail for preview
	AlbumID       primitive.ObjectID   `bson:"album_id"`              // Album reference
	UserID        primitive.ObjectID   `bson:"uploader_user_id"`      // Uploader reference
	Description   string               `bson:"description,omitempty"` // Optional description
	UploadedAt    time.Time            `bson:"uploaded_at"`           // Upload timestamp
	FacesID       []primitive.ObjectID `bson:"faces_id,omitempty"`    // List of recognized face IDs
	Width         int                  `bson:"width,omitempty"`       // Image width in pixels
	Height        int                  `bson:"height,omitempty"`      // Image height in pixels
	//FileSize    int64                `bson:"file_size,omitempty"`   // File size in bytes
}

// SmartFrame Represents a smart frame device
type SmartFrame struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty"`
	OwnerID      primitive.ObjectID   `bson:"owner_id" binding:"required"` // Owner's user ID
	GiftedByID   primitive.ObjectID   `bson:"gifted_by_id,omitempty"`      // Gifter's user ID
	CreatedAt    time.Time            `bson:"created_at"`                  // Device creation timestamp
	FirstBoot    time.Time            `bson:"first_boot"`                  // Timestamp of first boot
	LoadedAlbums []primitive.ObjectID `bson:"loaded_albums_id,omitempty"`  // Preloaded albums
}

// PictureData Represents the binary data of a picture
type PictureData struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Data []byte             `bson:"data,omitempty"`
}
