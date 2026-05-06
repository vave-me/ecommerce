#!/bin/sh
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE commondb TEMPLATE template0;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "commondb" <<-EOSQL
  -- Apply to keep modifications to the created_at column from being made
  CREATE OR REPLACE FUNCTION created_at_trigger()
  RETURNS TRIGGER AS \$\$
  BEGIN
    NEW.created_at := OLD.created_at;
    RETURN NEW;
  END;
  \$\$ language 'plpgsql';

  -- Apply to a table to automatically update update_at columns
  CREATE OR REPLACE FUNCTION updated_at_trigger()
  RETURNS TRIGGER AS \$\$
  BEGIN
     IF row(NEW.*) IS DISTINCT FROM row(OLD.*) THEN
        NEW.updated_at = NOW();
        RETURN NEW;
     ELSE
        RETURN OLD;
     END IF;
  END;
  \$\$ language 'plpgsql';
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE baskets TEMPLATE commondb;

  CREATE USER baskets_user WITH ENCRYPTED PASSWORD 'baskets_pass';
  GRANT USAGE ON SCHEMA public TO baskets_user;
  GRANT CREATE, CONNECT ON DATABASE baskets TO baskets_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "baskets" <<-EOSQL
  CREATE SCHEMA baskets;
  GRANT CREATE, USAGE ON SCHEMA baskets TO baskets_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE categories TEMPLATE commondb;

  CREATE USER categories_user WITH ENCRYPTED PASSWORD 'categories_pass';
  GRANT USAGE ON SCHEMA public TO categories_user;
  GRANT CREATE, CONNECT ON DATABASE categories TO categories_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "categories" <<-EOSQL
  CREATE SCHEMA categories;
  GRANT CREATE, USAGE ON SCHEMA categories TO categories_user;
EOSQL


psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE notifications TEMPLATE commondb;

  CREATE USER notifications_user WITH ENCRYPTED PASSWORD 'notifications_pass';
  GRANT USAGE ON SCHEMA public TO notifications_user;
  GRANT CREATE, CONNECT ON DATABASE notifications TO notifications_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "notifications" <<-EOSQL
  CREATE SCHEMA notifications;
  GRANT CREATE, USAGE ON SCHEMA notifications TO notifications_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE ordering TEMPLATE commondb;

  CREATE USER ordering_user WITH ENCRYPTED PASSWORD 'ordering_pass';
  GRANT USAGE ON SCHEMA public TO ordering_user;
  GRANT CREATE, CONNECT ON DATABASE ordering TO ordering_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "ordering" <<-EOSQL
  CREATE SCHEMA ordering;
  GRANT CREATE, USAGE ON SCHEMA ordering TO ordering_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE payments TEMPLATE commondb;

  CREATE USER payments_user WITH ENCRYPTED PASSWORD 'payments_pass';
  GRANT USAGE ON SCHEMA public TO payments_user;
  GRANT CREATE, CONNECT ON DATABASE payments TO payments_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "payments" <<-EOSQL
  CREATE SCHEMA payments;
  GRANT CREATE, USAGE ON SCHEMA payments TO payments_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE search TEMPLATE commondb;

  CREATE USER search_user WITH ENCRYPTED PASSWORD 'search_pass';
  GRANT USAGE ON SCHEMA public TO search_user;
  GRANT CREATE, CONNECT ON DATABASE search TO search_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "search" <<-EOSQL
  CREATE SCHEMA search;
  GRANT CREATE, USAGE ON SCHEMA search TO search_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE geocoding TEMPLATE commondb;

  CREATE USER geocoding_user WITH ENCRYPTED PASSWORD 'geocoding_pass';
  GRANT USAGE ON SCHEMA public TO geocoding_user;
  GRANT CREATE, CONNECT ON DATABASE geocoding TO geocoding_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "geocoding" <<-EOSQL
  CREATE SCHEMA geocoding;
  GRANT CREATE, USAGE ON SCHEMA geocoding TO geocoding_user;
EOSQL


psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE users TEMPLATE commondb;

  CREATE USER users_user WITH ENCRYPTED PASSWORD 'users_pass';
  GRANT USAGE ON SCHEMA public TO users_user;
  GRANT CREATE, CONNECT ON DATABASE users TO users_user;
EOSQL

# Now we connect as superuser to 'users' DB and enable PostGIS
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "users" <<-EOSQL
  -- Install PostGIS in the 'users' DB
  CREATE EXTENSION IF NOT EXISTS postgis;

  CREATE SCHEMA users;
  GRANT CREATE, USAGE ON SCHEMA users TO users_user;
EOSQL


psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE products TEMPLATE commondb;

  CREATE USER products_user WITH ENCRYPTED PASSWORD 'products_pass';
  GRANT USAGE ON SCHEMA public TO products_user;
  GRANT CREATE, CONNECT ON DATABASE products TO products_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "products" <<-EOSQL
  -- Install PostGIS in the 'products' DB
  CREATE EXTENSION IF NOT EXISTS postgis;

  CREATE SCHEMA products;
  GRANT CREATE, USAGE ON SCHEMA products TO products_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE posts TEMPLATE commondb;

  CREATE USER posts_user WITH ENCRYPTED PASSWORD 'posts_pass';
  GRANT USAGE ON SCHEMA public TO posts_user;
  GRANT CREATE, CONNECT ON DATABASE posts TO posts_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "posts" <<-EOSQL
  -- Install PostGIS in the 'posts' DB
  CREATE EXTENSION IF NOT EXISTS postgis;

  CREATE SCHEMA posts;
  GRANT CREATE, USAGE ON SCHEMA posts TO posts_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE deals TEMPLATE commondb;

  CREATE USER deals_user WITH ENCRYPTED PASSWORD 'deals_pass';
  GRANT USAGE ON SCHEMA public TO deals_user;
  GRANT CREATE, CONNECT ON DATABASE deals TO deals_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "deals" <<-EOSQL
  -- Install PostGIS in the 'deals' DB
  CREATE EXTENSION IF NOT EXISTS postgis;

  CREATE SCHEMA deals;
  GRANT CREATE, USAGE ON SCHEMA deals TO deals_user;
EOSQL


psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE messages TEMPLATE commondb;
  CREATE USER messages_user WITH ENCRYPTED PASSWORD 'messages_pass';
  GRANT USAGE ON SCHEMA public TO messages_user;
  GRANT CREATE, CONNECT ON DATABASE messages TO messages_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "messages" <<-EOSQL
  CREATE SCHEMA messages;
  GRANT CREATE, USAGE ON SCHEMA messages TO messages_user;
EOSQL


psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE comments TEMPLATE commondb;

  CREATE USER comments_user WITH ENCRYPTED PASSWORD 'comments_pass';
  GRANT USAGE ON SCHEMA public TO comments_user;
  GRANT CREATE, CONNECT ON DATABASE comments TO comments_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "comments" <<-EOSQL
  CREATE SCHEMA comments;
  GRANT CREATE, USAGE ON SCHEMA comments TO comments_user;
EOSQL


psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE activity TEMPLATE commondb;

  CREATE USER activity_user WITH ENCRYPTED PASSWORD 'activity_pass';
  GRANT USAGE ON SCHEMA public TO activity_user;
  GRANT CREATE, CONNECT ON DATABASE activity TO activity_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "activity" <<-EOSQL
  CREATE SCHEMA activity;
  GRANT CREATE, USAGE ON SCHEMA activity TO activity_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE support TEMPLATE commondb;

  CREATE USER support_user WITH ENCRYPTED PASSWORD 'support_pass';
  GRANT USAGE ON SCHEMA public TO support_user;
  GRANT CREATE, CONNECT ON DATABASE support TO support_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "support" <<-EOSQL
  CREATE SCHEMA support;
  GRANT CREATE, USAGE ON SCHEMA support TO support_user;
EOSQL


psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE newsletters TEMPLATE commondb;

  CREATE USER newsletters_user WITH ENCRYPTED PASSWORD 'newsletters_pass';
  GRANT USAGE ON SCHEMA public TO newsletters_user;
  GRANT CREATE, CONNECT ON DATABASE newsletters TO newsletters_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "newsletters" <<-EOSQL
  CREATE SCHEMA newsletters;
  GRANT CREATE, USAGE ON SCHEMA newsletters TO newsletters_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE offers TEMPLATE commondb;

  CREATE USER offers_user WITH ENCRYPTED PASSWORD 'offers_pass';
  GRANT USAGE ON SCHEMA public TO offers_user;
  GRANT CREATE, CONNECT ON DATABASE offers TO offers_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "offers" <<-EOSQL
  CREATE SCHEMA offers;
  GRANT CREATE, USAGE ON SCHEMA offers TO offers_user;
EOSQL


psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE media TEMPLATE commondb;

  CREATE USER media_user WITH ENCRYPTED PASSWORD 'media_pass';
  GRANT USAGE ON SCHEMA public TO media_user;
  GRANT CREATE, CONNECT ON DATABASE media TO media_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "media" <<-EOSQL
  CREATE SCHEMA media;
  GRANT CREATE, USAGE ON SCHEMA media TO media_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE mailer TEMPLATE commondb;

  CREATE USER mailer_user WITH ENCRYPTED PASSWORD 'mailer_pass';
  GRANT USAGE ON SCHEMA public TO mailer_user;
  GRANT CREATE, CONNECT ON DATABASE mailer TO mailer_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "mailer" <<-EOSQL
  CREATE SCHEMA mailer;
  GRANT CREATE, USAGE ON SCHEMA mailer TO mailer_user;
EOSQL


psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE wishlists TEMPLATE commondb;

  CREATE USER wishlists_user WITH ENCRYPTED PASSWORD 'wishlists_pass';
  GRANT USAGE ON SCHEMA public TO wishlists_user;
  GRANT CREATE, CONNECT ON DATABASE wishlists TO wishlists_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "wishlists" <<-EOSQL
  CREATE SCHEMA wishlists;
  GRANT CREATE, USAGE ON SCHEMA wishlists TO wishlists_user;
EOSQL


psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE shipping TEMPLATE commondb;
  CREATE USER shipping_user WITH ENCRYPTED PASSWORD 'shipping_pass';
  GRANT USAGE ON SCHEMA public TO shipping_user;
  GRANT CREATE, CONNECT ON DATABASE shipping TO shipping_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "shipping" <<-EOSQL
  CREATE SCHEMA shipping;
  GRANT CREATE, USAGE ON SCHEMA shipping TO shipping_user;
EOSQL


psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE following TEMPLATE commondb;

  CREATE USER following_user WITH ENCRYPTED PASSWORD 'following_pass';
  GRANT USAGE ON SCHEMA public TO following_user;
  GRANT CREATE, CONNECT ON DATABASE following TO following_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "following" <<-EOSQL
  CREATE SCHEMA following;
  GRANT CREATE, USAGE ON SCHEMA following TO following_user;
EOSQL


psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE reviews TEMPLATE commondb;

  CREATE USER reviews_user WITH ENCRYPTED PASSWORD 'reviews_pass';
  GRANT USAGE ON SCHEMA public TO reviews_user;
  GRANT CREATE, CONNECT ON DATABASE reviews TO reviews_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "reviews" <<-EOSQL
  CREATE SCHEMA reviews;
  GRANT CREATE, USAGE ON SCHEMA reviews TO reviews_user;
EOSQL

# metrics
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE metrics TEMPLATE commondb;

  CREATE USER metrics_user WITH ENCRYPTED PASSWORD 'metrics_pass';
  GRANT USAGE ON SCHEMA public TO metrics_user;
  GRANT CREATE, CONNECT ON DATABASE metrics TO metrics_user;
EOSQL

# Now we connect as superuser to 'users' DB and enable PostGIS
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "metrics" <<-EOSQL
  -- Install PostGIS in the 'metrics' DB
  CREATE EXTENSION IF NOT EXISTS postgis;

  CREATE SCHEMA metrics;
  GRANT CREATE, USAGE ON SCHEMA metrics TO metrics_user;
EOSQL

#assistans
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE assistants TEMPLATE commondb;

  CREATE USER assistants_user WITH ENCRYPTED PASSWORD 'assistants_pass';
  GRANT USAGE ON SCHEMA public TO assistants_user;
  GRANT CREATE, CONNECT ON DATABASE assistants TO assistants_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "assistants" <<-EOSQL
  CREATE SCHEMA assistants;
  GRANT CREATE, USAGE ON SCHEMA assistants TO assistants_user;
EOSQL


psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE cosec TEMPLATE commondb;

  CREATE USER cosec_user WITH ENCRYPTED PASSWORD 'cosec_pass';
  GRANT USAGE ON SCHEMA public TO cosec_user;
  GRANT CREATE, CONNECT ON DATABASE cosec TO cosec_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "cosec" <<-EOSQL
  CREATE SCHEMA cosec;
  GRANT CREATE, USAGE ON SCHEMA cosec TO cosec_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE services TEMPLATE commondb;

  CREATE USER services_user WITH ENCRYPTED PASSWORD 'services_pass';
  GRANT USAGE ON SCHEMA public TO services_user;
  GRANT CREATE, CONNECT ON DATABASE services TO services_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "services" <<-EOSQL
  -- Install PostGIS in the 'services' DB
  CREATE EXTENSION IF NOT EXISTS postgis;

  CREATE SCHEMA services;
  GRANT CREATE, USAGE ON SCHEMA services TO services_user;
EOSQL


psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE scheduler TEMPLATE commondb;

  CREATE USER scheduler_user WITH ENCRYPTED PASSWORD 'scheduler_pass';
  GRANT USAGE ON SCHEMA public TO scheduler_user;
  GRANT CREATE, CONNECT ON DATABASE scheduler TO scheduler_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "scheduler" <<-EOSQL
  CREATE SCHEMA scheduler;
  GRANT CREATE, USAGE ON SCHEMA scheduler TO scheduler_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE erp TEMPLATE commondb;

  CREATE USER erp_user WITH ENCRYPTED PASSWORD 'erp_pass';
  GRANT USAGE ON SCHEMA public TO erp_user;
  GRANT CREATE, CONNECT ON DATABASE erp TO erp_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "erp" <<-EOSQL
  CREATE SCHEMA erp;
  GRANT CREATE, USAGE ON SCHEMA erp TO erp_user;
EOSQL


#managers
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE managers TEMPLATE commondb;

  CREATE USER managers_user WITH ENCRYPTED PASSWORD 'managers_pass';
  GRANT USAGE ON SCHEMA public TO managers_user;
  GRANT CREATE, CONNECT ON DATABASE managers TO managers_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "managers" <<-EOSQL
  CREATE SCHEMA managers;
  GRANT CREATE, USAGE ON SCHEMA managers TO managers_user;
EOSQL

#managers
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE merchant TEMPLATE commondb;

  CREATE USER merchant_user WITH ENCRYPTED PASSWORD 'merchant_pass';
  GRANT USAGE ON SCHEMA public TO managers_user;
  GRANT CREATE, CONNECT ON DATABASE merchant TO merchant_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "merchant" <<-EOSQL
  CREATE SCHEMA merchant;
  GRANT CREATE, USAGE ON SCHEMA merchant TO merchant_user;
EOSQL