# Task-2 

Based on the scenario and requirement mentioned, the following REST API endpoints can be implemented:

1.Endpoint URL: /apples/lots 
HTTP Method: POST 
Request Payload:

{
    "cultivar": "Red Dacca",
    "country_of_origin": "Costa Rica",
    "harvesting_date": "2018-07-27",
    "total_weight": 500
}
Response Payload:

{
    "lot_id": "abcd-1234",
    "status": "created"
}
Response Status Code: 201 (Created)

Description: This endpoint is used to create a new lot of apples for sale. The request payload contains the details of the lot, including cultivar, country of origin, harvesting date, and total weight. If the weight of the lot is less than the minimum allowed weight (i.e. 1000 kg), an error response will be returned with status code 400 (Bad Request). If the creation of the lot is successful, a unique lot_id will be returned in the response payload.

2. Endpoint URL: /apples/lots/{lot_id} 
HTTP Method: PUT 
Request Payload:

{
    "harvesting_date": "2018-06-14"
}
Response Payload:

{
    "status": "updated"
}
Response Status Code: 200 (OK)

Description: This endpoint is used to update the harvesting date of an existing lot of apples. The lot_id parameter in the URL identifies the lot to be updated. The request payload contains the new harvesting date. If the update is successful, a success response will be returned with status code 200 (OK).

3. Endpoint URL: /apples/lots/{lot_id}/auction 
HTTP Method: POST 
Request Payload:

{
    "start_date": "2022-04-01T12:00:00Z",
    "end_date": "2022-04-02T12:00:00Z",
    "starting_bid_price": 10.00
}
Response Payload:

{
    "auction_id": "abcd-1234",
    "status": "started"
}
Response Status Code: 201 (Created)

Description: This endpoint is used to start an auction.
