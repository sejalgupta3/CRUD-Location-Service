CRUD Location Service
======
The location service has the following REST endpoints to store and retrieve locations.
For lookup coordinates of a location, Google Map Api has been used.

Example:
Get coordinates of 1600 Amphitheatre Parkway, Mountain View, CA.
http://maps.google.com/maps/api/geocode/json?address=1600+Amphitheatre+Parkway,+Mountain+View,+CA&sensor=false

* Create New Location - POST        /locations
```
Request
{
   "name" : "John Smith",
   "address" : "123 Main St",
   "city" : "San Francisco",
   "state" : "CA",
   "zip" : "94113"
}
```
```
Response
{
   "id" : 12345,
   "name" : "John Smith",
   "address" : "123 Main St",
   "city" : "San Francisco",
   "state" : "CA",
   "zip" : "94113",
   "coordinate" : { 
      "lat" : 38.4220352,
     "lng" : -222.0841244
   }
}

```
* Get a Location - GET        /locations/{location_id}
```
Request:
GET /locations/12345
```
```
Response:
{
   "id" : 12345,
   "name" : "John Smith",
   "address" : "123 Main St",
   "city" : "San Francisco",
   "state" : "CA",
   "zip" : "94113",
   "coordinate" : { 
      "lat" : 38.4220352,
     "lng" : -222.0841244
   }
}

```
* Update a Location - PUT /locations/{location_id}
```
Request:
{
   "address" : "1600 Amphitheatre Parkway",
   "city" : "Mountain View",
   "state" : "CA",
   "zip" : "94043"
}
```
```
Response:
{
   "id" : 12345,
   "name" : "John Smith",
   "address" : "1600 Amphitheatre Parkway",
   "city" : "Mountain View",
   "state" : "CA",
   "zip" : "94043"
   "coordinate" : { 
      "lat" : 37.4220352,
     "lng" : -122.0841244
   }
}
```
* Delete a Location - DELETE /locations/{location_id}
```
Request:
DELETE  /locations/12345
```
```
Response:
HTTP Response Code: 200
```
