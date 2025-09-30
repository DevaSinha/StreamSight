FROM gocv/opencv:4.12.0

WORKDIR /app

# Copy go.mod and go.sum to download deps first
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Copy Haar cascade for face detection
COPY haarcascade_frontalface_default.xml .

# Build the application
RUN go build -o workerbin .

CMD ["./workerbin"]
