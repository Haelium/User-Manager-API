# User Management API

This was a take-home assignment for my application for a previous job. A hard requirement was using redis as the backing datastore. Another hard requirement was deploying the system on kubernetes and making it horizontally scalable, for this reason the service is stateless.

All interactions with the database happen through an interface called DatabaseInterface, which RedisHashConn implements. For this reason, it is trivial to replace redis with another database in this codebase.
