# PHP-FPM Docker Image for 20i-style hosting
ARG PHP_VERSION=8.4
FROM php:${PHP_VERSION}-fpm-alpine

# Install common PHP extensions found in shared hosting
RUN apk add --no-cache \
    freetype-dev \
    libjpeg-turbo-dev \
    libpng-dev \
    libwebp-dev \
    libzip-dev \
    zip \
    unzip \
    curl-dev \
    oniguruma-dev \
    autoconf \
    g++ \
    make \
    && docker-php-ext-configure gd --with-freetype --with-jpeg --with-webp \
    && docker-php-ext-install -j$(nproc) \
        gd \
        mysqli \
        pdo_mysql \
        zip \
        curl \
        mbstring \
        exif \
        bcmath \
        opcache

# Note: fileinfo and iconv are already built into PHP 8.4
# No need to install them separately

# Copy custom PHP configuration
COPY php.ini /usr/local/etc/php/conf.d/custom.ini

# Set working directory
WORKDIR /var/www/site

# Configure PHP-FPM to listen on port 9000
RUN echo "listen = 0.0.0.0:9000" >> /usr/local/etc/php-fpm.d/www.conf \
    && echo "pm.max_children = 20" >> /usr/local/etc/php-fpm.d/www.conf \
    && echo "pm.start_servers = 3" >> /usr/local/etc/php-fpm.d/www.conf \
    && echo "pm.min_spare_servers = 2" >> /usr/local/etc/php-fpm.d/www.conf \
    && echo "pm.max_spare_servers = 5" >> /usr/local/etc/php-fpm.d/www.conf

# Expose PHP-FPM port
EXPOSE 9000

# Start PHP-FPM
CMD ["php-fpm"]
