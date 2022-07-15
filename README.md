# Zappr

Zappr is a low code, minimalistic API server for all kinds of applications. Built with Go for blazing speed, Zappr can get your backend up and running with just a simple docker command (which also means it is self deployable)

# Getting Started

Getting started with Zappr is as simple as boiling eggs. All you need to do is use the docker run command which will pull the latest Zappr image from Docker Hub, mix it with some essential environment variables, and run the image locally. It is all self deployable. Alternatively, one can also clone the repository, and build the Docker image manually.

Below are some important entities to know about in Zappr

## Tenants

The root entity of everything that is Zappr. In order for Zappr's APIs to be used, a tenant needs to be created using Zappr's APIs, and every other operation can take place within that tenant. Of course, that means you can manage multiple tenants in one deployment. I can smell corporate people already :)

## Users

Every tenant needs an authenticated and authorized user to do operations on application-specific data. You first sign up an user, and the subsequently log that user in, to get a JWT. That JWT would be needed to orchestrate all application specific data

## Roles

Every tenant can have roles that can be assigned to aforementioned users. By default, o signup, every user is assigned the Normal User role, which only has read permissions. APIs can be used to elevate an user to higher permissions

## User Metadata

The secret sauce that makes Zappr awesome is this. User metadata is nothing but application specific data. The extensive API suite caters to all kinds of CRUD needs that an application might need.

# Features in the pipeline

Right now, Zappr is in alpha, as in, only the most rudimentary of all entities are included in its API suite. It is only sane that Zappr will evolve further.

## SMTP APIs

APIs to programmatically fire emails

## Static Storage APIs

API suite to address storage of images, videos or any other type of static files

## External Connectors

External connectors would allow Zappr to communicate to external services.

## Webhooks

Zappr would have a full fledged API suite for webhooks

# API Documentation

This Postman collection is as extensive as it gets when it comes to the API documentation
