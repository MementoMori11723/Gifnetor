# Stage 1: Build Stage with python:3.10-alpine for a smaller footprint
FROM python:3.10-alpine AS builder

# Set environment variables to avoid Python output buffering
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1

# Install build dependencies and required libraries
RUN apk add --no-cache \
    build-base \
    gcc \
    libc-dev \
    libjpeg-turbo-dev \
    zlib-dev \
    python3-dev \
    musl-dev \
    bash \
    jpeg-dev \
    freetype-dev \
    pkgconfig \
    py3-pip \
    py3-setuptools \
    py3-wheel \
    openblas-dev

# Set the working directory in the container
WORKDIR /app

# Copy requirements file and install dependencies
COPY requirements.txt .

# Install Python dependencies without caching to reduce the final image size
RUN pip install --no-cache-dir --user -r requirements.txt

# Stage 2: Final Stage - Use a minimal Alpine image for runtime
FROM python:3.10-alpine

# Install runtime dependencies only
RUN apk add --no-cache \
    libjpeg-turbo \
    zlib \
    libstdc++ \
    jpeg-dev \
    bash \
    openblas

# Set environment variables
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1

# Set the working directory in the container
WORKDIR /app

# Copy installed dependencies from the builder stage
COPY --from=builder /root/.local /root/.local

# Ensure the Python scripts in the user folder are in the PATH
ENV PATH=/root/.local/bin:$PATH

# Copy the application code
COPY . .

# Expose port for streamlit
EXPOSE 8501

# Run the app
CMD ["streamlit", "run", "app.py"]


