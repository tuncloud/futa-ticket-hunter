Bạn cần sử dụng các headers sau đây khi thực hiện gọi các API bên dưới
```
User-Agent: FUTA/7.36.4 (com.client.facecar; build:1; iOS 26.2.0) Alamofire/5.9.1
X-App-Version: 7.36.4
X-Access-Token: eyJhbGciOiJSUzI1NiIsImtpZCI6IjM3MzAwNzY5YTA3ZTA1MTE2ZjdlNTEzOGZhOTA5MzY4NWVlYmMyNDAiLCJ0eXAiOiJKV1QifQ.eyJmdWxsX25hbWUiOiLEkOG7lyBRdeG7kWMgVHXhuqVuIiwiY3VzdG9tX3VpZCI6MjE3MDE2NywiYXV0aG9yaXRpZXMiOnsidXNlciI6e319LCJhZG1pbl96b25lX2lkIjo2NSwiaXNzIjoiaHR0cHM6Ly9zZWN1cmV0b2tlbi5nb29nbGUuY29tL2ZhY2VjYXItMjlhZTciLCJhdWQiOiJmYWNlY2FyLTI5YWU3IiwiYXV0aF90aW1lIjoxNzMwNTQ1NDYyLCJ1c2VyX2lkIjoiY3FhaFE0cVhreE44cVBIUm8yUVJlZFhXaGhPMiIsInN1YiI6ImNxYWhRNHFYa3hOOHFQSFJvMlFSZWRYV2hoTzIiLCJpYXQiOjE3NzQ4ODE5NDcsImV4cCI6MTc3NDg4NTU0NywicGhvbmVfbnVtYmVyIjoiKzg0MzY3NzE3NzE0IiwiZmlyZWJhc2UiOnsiaWRlbnRpdGllcyI6eyJwaG9uZSI6WyIrODQzNjc3MTc3MTQiXX0sInNpZ25faW5fcHJvdmlkZXIiOiJjdXN0b20ifX0.YtpP8v10UCwpssqpeypSrv-KgfLYYPDgroq_IPHpw8kE62UPTGf_sEFfBNbvUiEOVRq8FlpreWdFsmpljCqJV_hRCXL0tOQaZcGPHzw3JptcliIqRgWr3DpJwAk6vqCCf35kVY6a7ybCph38_dj1QwrnO8C3pv5oOmNsgel_4U7OWqcpoOVIYzP81JxeXLKPjZi0m0q46DXLKEjZ-lv2Pr9G5WLN5VCsgHSrOJr0MeKgwxplmfqSTHsp2gVgTnYF0YSZys1xCQiEAaVggvVPveQjwb_1dk_zFQy8Bnc8mwuKE0L-IRKCN-2aptjwlDnfTrIESgwRkUnGRwIu7xoKJQ
X-Channel: mobile_app
```

Trong đó giá trị của header X-Access-Token được lấy bằng cách:
- fetch url https://futabus.vn
- extract token từ trong raw response trả về, có 1 field tên là token trong raw html response

Search pickup points
```
https://api-online.futabus.vn/vato/v1/search/pickup-point?keyword=ho chi minh&page=0&size=50
```

response
```
{
    "requestId": "019d3f47-d85b-7b8a-a7d3-db22bd0a98f3",
    "status": 200,
    "error": null,
    "data": {
        "@type": "type.googleapis.com/vn.futa.vato.buslines.api.PaginationData",
        "page": 0,
        "size": 50,
        "total": 1000,
        "items": [
            {
                "@type": "type.googleapis.com/vn.futa.vato.buslines.api.search.v1.GroupMetadataPickupPointSearch",
                "districtId": "48c5de67-fcce-1925-8772-11912f57446d",
                "districtName": "Quận 5",
                "provinceName": "TP. Hồ Chí Minh",
                "group": [
                    {
                        "departmentId": "11f09e78-0e69-3b42-a90c-42964ca74856",
                        "departmentName": "231-233 Lê Hồng Phong",
                        "departmentAddress": "231 Đ. Lê Hồng Phong, Phường 4, Quận 5, Thành phố Hồ Chí Minh, Vietnam",
                        "departmentTime": -15,
                        "areaId": "41146f8e-0fa0-d3f3-90e3-23ef2db2bf81",
                        "provinceId": "421741fc-2746-7e55-a5a4-eec23c50c986",
                        "provinceName": "TP. Hồ Chí Minh",
                        "districtId": "48c5de67-fcce-1925-8772-11912f57446d",
                        "districtName": "Quận 5",
                        "type": 0,
                        "latitude": 10.822683,
                        "longitude": 106.62927
                    }
                ]
            },
            {
                "@type": "type.googleapis.com/vn.futa.vato.buslines.api.search.v1.GroupMetadataPickupPointSearch",
                "districtId": "449e8b60-5f01-c5af-93c5-d88015f41834",
                "districtName": "Quận 8",
                "provinceName": "TP. Hồ Chí Minh",
                "group": [
                    {
                        "departmentId": "11f09e78-0e69-45dc-a90c-42964ca74856",
                        "departmentName": "43 Nguyễn Cư Trinh",
                        "departmentAddress": "43 Đ. Nguyễn Cư Trinh, Phường Nguyễn Cư Trinh, Quận 1, Thành phố Hồ Chí Minh, Vietnam",
                        "departmentTime": -15,
                        "areaId": "4e2c2f9a-3bc4-b8d6-a0b4-0110810e195c",
                        "provinceId": "421741fc-2746-7e55-a5a4-eec23c50c986",
                        "provinceName": "TP. Hồ Chí Minh",
                        "districtId": "449e8b60-5f01-c5af-93c5-d88015f41834",
                        "districtName": "Quận 8",
                        "type": 0,
                        "latitude": 0,
                        "longitude": 0
                    }
                ]
            },
            {
                "@type": "type.googleapis.com/vn.futa.vato.buslines.api.search.v1.GroupMetadataPickupPointSearch",
                "districtId": "4018efc9-ab9a-ee82-9c74-af0258fd78ad",
                "districtName": "Quận 10",
                "provinceName": "TP. Hồ Chí Minh",
                "group": [
                    {
                        "departmentId": "747e14ea-409b-4a54-a2d0-da18915969f1",
                        "departmentName": "VP Thành Thái",
                        "departmentAddress": "108 Đ. Thành Thái, Phường 14, Quận 10, Thành phố Hồ Chí Minh, Vietnam",
                        "departmentTime": -45,
                        "areaId": "4c556284-b898-9405-9e5a-746f9ef3436a",
                        "provinceId": "421741fc-2746-7e55-a5a4-eec23c50c986",
                        "provinceName": "TP. Hồ Chí Minh",
                        "districtId": "4018efc9-ab9a-ee82-9c74-af0258fd78ad",
                        "districtName": "Quận 10",
                        "type": 0,
                        "latitude": 10.774236,
                        "longitude": 106.66468
                    }
                ]
            },
            {
                "@type": "type.googleapis.com/vn.futa.vato.buslines.api.search.v1.GroupMetadataPickupPointSearch",
                "districtId": "4a3d3753-7e7e-d218-b44a-bddf360c540d",
                "districtName": "Quận 4",
                "provinceName": "TP. Hồ Chí Minh",
                "group": [
                    {
                        "departmentId": "11f09e78-0e69-3b7a-a90c-42964ca74856",
                        "departmentName": "Xa Lộ Hà Nội",
                        "departmentAddress": "798 Song Hành Xa Lộ Hà Nội, Hiệp Phú, Quận 9, Thành phố Hồ Chí Minh, Vietnam",
                        "departmentTime": 30,
                        "areaId": "47142f67-ab42-3ba6-8b8e-3b18a33d8956",
                        "provinceId": "421741fc-2746-7e55-a5a4-eec23c50c986",
                        "provinceName": "TP. Hồ Chí Minh",
                        "districtId": "4a3d3753-7e7e-d218-b44a-bddf360c540d",
                        "districtName": "Quận 4",
                        "type": 0,
                        "latitude": 0,
                        "longitude": 0
                    }
                ]
            },
            {
                "@type": "type.googleapis.com/vn.futa.vato.buslines.api.search.v1.GroupMetadataPickupPointSearch",
                "districtId": "47138b71-3e58-8da4-a951-29fff3344826",
                "districtName": "Bình Thạnh",
                "provinceName": "TP. Hồ Chí Minh",
                "group": [
                    {
                        "departmentId": "b1fa5c43-e2f3-458e-87a7-cb39c5fa2afc",
                        "departmentName": "Bệnh Viện Ung Bướu Bình Thạnh",
                        "departmentAddress": "3 Nơ Trang Long, Phường 7, Bình Thạnh, Thành phố Hồ Chí Minh, Vietnam",
                        "departmentTime": -90,
                        "areaId": "48daeb8f-6baa-e944-9975-97dfd2b0ced3",
                        "provinceId": "421741fc-2746-7e55-a5a4-eec23c50c986",
                        "provinceName": "TP. Hồ Chí Minh",
                        "districtId": "47138b71-3e58-8da4-a951-29fff3344826",
                        "districtName": "Bình Thạnh",
                        "type": 0,
                        "latitude": 10.804886,
                        "longitude": 106.69462
                    },
                    {
                        "departmentId": "11f09e78-0e69-3b89-a90c-42964ca74856",
                        "departmentName": "Bến xe Miền Đông Cũ",
                        "departmentAddress": "292 Đinh Bộ Lĩnh, P.26, Q.Bình Thạnh, TP HCM",
                        "departmentTime": 30,
                        "areaId": "44f439ff-cdf0-e0a5-bed4-b68f28835b2d",
                        "provinceId": "421741fc-2746-7e55-a5a4-eec23c50c986",
                        "provinceName": "TP. Hồ Chí Minh",
                        "districtId": "47138b71-3e58-8da4-a951-29fff3344826",
                        "districtName": "Bình Thạnh",
                        "type": 0,
                        "latitude": 0,
                        "longitude": 0
                    }
                ]
            },
            {
                "@type": "type.googleapis.com/vn.futa.vato.buslines.api.search.v1.GroupMetadataPickupPointSearch",
                "districtId": "4f524fb5-2f4a-1aab-a27f-0b6fc224420d",
                "districtName": "Hóc Môn",
                "provinceName": "TP. Hồ Chí Minh",
                "group": [
                    {
                        "departmentId": "11f09e78-0e69-425c-a90c-42964ca74856",
                        "departmentName": "Bến Xe An Sương",
                        "departmentAddress": "Quốc Lộ 22, Ấp Đông Lân, Bà Điểm, Hóc Môn, TP Hồ Chí Minh",
                        "departmentTime": 60,
                        "areaId": "43005da7-ae5b-837b-b115-130bef1acc95",
                        "provinceId": "421741fc-2746-7e55-a5a4-eec23c50c986",
                        "provinceName": "TP. Hồ Chí Minh",
                        "districtId": "4f524fb5-2f4a-1aab-a27f-0b6fc224420d",
                        "districtName": "Hóc Môn",
                        "type": 0,
                        "latitude": 10.843817,
                        "longitude": 106.613785
                    }
                ]
            },
            {
                "@type": "type.googleapis.com/vn.futa.vato.buslines.api.search.v1.GroupMetadataPickupPointSearch",
                "districtId": "4bbb343c-0ce4-c4d9-86ac-d1721fb31c47",
                "districtName": "Quận 9",
                "provinceName": "TP. Hồ Chí Minh",
                "group": [
                    {
                        "departmentId": "11f09e78-0e69-4706-a90c-42964ca74856",
                        "departmentName": "Bến xe Miền Đông Mới",
                        "departmentAddress": "501 Hoàng Hữu Nam, Long Bình, Thủ Đức, Thành phố Hồ Chí Minh, Vietnam",
                        "departmentTime": 0,
                        "areaId": "48201c1f-df53-aa10-b243-608c7f4a9548",
                        "provinceId": "421741fc-2746-7e55-a5a4-eec23c50c986",
                        "provinceName": "TP. Hồ Chí Minh",
                        "districtId": "4bbb343c-0ce4-c4d9-86ac-d1721fb31c47",
                        "districtName": "Quận 9",
                        "type": 0,
                        "latitude": 10.878811,
                        "longitude": 106.8166
                    }
                ]
            },
            {
                "@type": "type.googleapis.com/vn.futa.vato.buslines.api.search.v1.GroupMetadataPickupPointSearch",
                "districtId": "4885ad91-cc5b-f727-8d20-869321307dd2",
                "districtName": "Bình Tân",
                "provinceName": "TP. Hồ Chí Minh",
                "group": [
                    {
                        "departmentId": "11f09e78-0e69-3b72-a90c-42964ca74856",
                        "departmentName": "Bến xe Miền Tây",
                        "departmentAddress": "395 Kinh Dương Vương , P.An Lạc , Q.Bình Tân , TP.HCM",
                        "departmentTime": 0,
                        "areaId": "46f6450f-e2dc-2112-9979-4e850cd42422",
                        "provinceId": "421741fc-2746-7e55-a5a4-eec23c50c986",
                        "provinceName": "TP. Hồ Chí Minh",
                        "districtId": "4885ad91-cc5b-f727-8d20-869321307dd2",
                        "districtName": "Bình Tân",
                        "type": 0,
                        "latitude": 10.822683,
                        "longitude": 106.62927
                    }
                ]
            },
            {
                "@type": "type.googleapis.com/vn.futa.vato.buslines.api.search.v1.GroupMetadataPickupPointSearch",
                "districtId": "408a6203-b330-fa14-9067-b59c35543c7a",
                "districtName": "Phan Thiết",
                "provinceName": "Bình Thuận",
                "group": [
                    {
                        "departmentId": "98972dd1-eb50-4260-a170-3ee3ea539fa9",
                        "departmentName": "VP Online",
                        "departmentAddress": "102 Đ. Trần Hưng Đạo, Phường Phạm Ngũ Lão, Quận 1, Thành phố Hồ Chí Minh, Vietnam",
                        "departmentTime": 300,
                        "areaId": "4361d1d3-9c43-e821-9cb1-62995506ef2a",
                        "provinceId": "4509ec3f-cf5f-a141-bba7-866864357aec",
                        "provinceName": "Bình Thuận",
                        "districtId": "408a6203-b330-fa14-9067-b59c35543c7a",
                        "districtName": "Phan Thiết",
                        "type": 0,
                        "latitude": 10.9504385,
                        "longitude": 108.110985
                    }
                ]
            },
            {
                "@type": "type.googleapis.com/vn.futa.vato.buslines.api.search.v1.GroupMetadataPickupPointSearch",
                "districtId": "44e5958e-8ab7-fa1f-94c8-201687589a83",
                "districtName": "Bắc Bình",
                "provinceName": "Bình Thuận",
                "group": [
                    {
                        "departmentId": "52f72df4-961b-46d3-a93e-bd9d95b9c07c",
                        "departmentName": "Bệnh viện Ung Bứu Quận 9",
                        "departmentAddress": "12 Đ. Số 400, Long Thạnh Mỹ, Thủ Đức, Thành phố Hồ Chí Minh, Vietnam",
                        "departmentTime": 300,
                        "areaId": "4d780db4-2f03-99e1-bb2b-5058578813f5",
                        "provinceId": "4509ec3f-cf5f-a141-bba7-866864357aec",
                        "provinceName": "Bình Thuận",
                        "districtId": "44e5958e-8ab7-fa1f-94c8-201687589a83",
                        "districtName": "Bắc Bình",
                        "type": 0,
                        "latitude": 10.871213,
                        "longitude": 106.80937
                    }
                ]
            },
            {
                "@type": "type.googleapis.com/vn.futa.vato.buslines.api.search.v1.GroupMetadataPickupPointSearch",
                "districtId": "46d8ee64-c59b-495a-b016-4f1a95714e3d",
                "districtName": "Quận 12",
                "provinceName": "TP. Hồ Chí Minh",
                "group": [
                    {
                        "departmentId": "11f09e78-0e69-475a-a90c-42964ca74856",
                        "departmentName": "Bến xe Ga",
                        "departmentAddress": "Bến xe Miền Đông mới, Bình Thắng, Dĩ An, Bình Dương, Việt Nam",
                        "departmentTime": 60,
                        "areaId": "4a5f7ae6-2cde-45b7-b40a-3da53a548fa2",
                        "provinceId": "421741fc-2746-7e55-a5a4-eec23c50c986",
                        "provinceName": "TP. Hồ Chí Minh",
                        "districtId": "46d8ee64-c59b-495a-b016-4f1a95714e3d",
                        "districtName": "Quận 12",
                        "type": 0,
                        "latitude": 10.880404,
                        "longitude": 106.814865
                    }
                ]
            }
        ],
        "others": [
            {
                "@type": "type.googleapis.com/vn.futa.vato.buslines.api.search.v1.Area",
                "id": "421741fc-2746-7e55-a5a4-eec23c50c986",
                "name": "TP. Hồ Chí Minh",
                "fullAddress": "TP. Hồ Chí Minh",
                "level": 2,
                "code": "TPHCM",
                "parentId": "4106b053-0d60-66a3-b5fd-1dc054e542ac",
                "tags": "TP. Ho Chi Minh",
                "formattedAddress": ""
            },
            {
                "@type": "type.googleapis.com/vn.futa.vato.buslines.api.search.v1.Area",
                "id": "4d9a969a-d469-d77b-b789-aacb9c393d81",
                "name": "Quận 1",
                "fullAddress": "Quận 1, TP. Hồ Chí Minh",
                "level": 3,
                "code": "",
                "parentId": "421741fc-2746-7e55-a5a4-eec23c50c986",
                "tags": "Quan 1, TP. Ho Chi Minh",
                "formattedAddress": ""
            },
            {
                "@type": "type.googleapis.com/vn.futa.vato.buslines.api.search.v1.Area",
                "id": "48d092d8-40a1-a4bf-a58e-c2ee62d8c2c0",
                "name": "Quận 2",
                "fullAddress": "Quận 2, TP. Hồ Chí Minh",
                "level": 3,
                "code": "",
                "parentId": "421741fc-2746-7e55-a5a4-eec23c50c986",
                "tags": "Quan 2, TP. Ho Chi Minh",
                "formattedAddress": ""
            },
            {
                "@type": "type.googleapis.com/vn.futa.vato.buslines.api.search.v1.Area",
                "id": "4dbe44a8-2645-bfc1-a74e-1f6fb3b6fd0f",
                "name": "Quận 3",
                "fullAddress": "Quận 3, TP. Hồ Chí Minh",
                "level": 3,
                "code": "",
                "parentId": "421741fc-2746-7e55-a5a4-eec23c50c986",
                "tags": "Quan 3, TP. Ho Chi Minh",
                "formattedAddress": ""
            },
            {
                "@type": "type.googleapis.com/vn.futa.vato.buslines.api.search.v1.Area",
                "id": "48c5de67-fcce-1925-8772-11912f57446d",
                "name": "Quận 5",
                "fullAddress": "Quận 5, TP. Hồ Chí Minh",
                "level": 3,
                "code": "",
                "parentId": "421741fc-2746-7e55-a5a4-eec23c50c986",
                "tags": "Quan 5, TP. Ho Chi Minh",
                "formattedAddress": ""
            },
            {
                "@type": "type.googleapis.com/vn.futa.vato.buslines.api.search.v1.Area",
                "id": "449e8b60-5f01-c5af-93c5-d88015f41834",
                "name": "Quận 8",
                "fullAddress": "Quận 8, TP. Hồ Chí Minh",
                "level": 3,
                "code": "",
                "parentId": "421741fc-2746-7e55-a5a4-eec23c50c986",
                "tags": "Quan 8, TP. Ho Chi Minh",
                "formattedAddress": ""
            },
            {
                "@type": "type.googleapis.com/vn.futa.vato.buslines.api.search.v1.Area",
                "id": "4bbb343c-0ce4-c4d9-86ac-d1721fb31c47",
                "name": "Quận 9",
                "fullAddress": "Quận 9, TP. Hồ Chí Minh",
                "level": 3,
                "code": "",
                "parentId": "421741fc-2746-7e55-a5a4-eec23c50c986",
                "tags": "Quan 9, TP. Ho Chi Minh",
                "formattedAddress": ""
            },
            {
                "@type": "type.googleapis.com/vn.futa.vato.buslines.api.search.v1.Area",
                "id": "4a3d3753-7e7e-d218-b44a-bddf360c540d",
                "name": "Quận 4",
                "fullAddress": "Quận 4, TP. Hồ Chí Minh",
                "level": 3,
                "code": "",
                "parentId": "421741fc-2746-7e55-a5a4-eec23c50c986",
                "tags": "Quan 4, TP. Ho Chi Minh",
                "formattedAddress": ""
            },
            {
                "@type": "type.googleapis.com/vn.futa.vato.buslines.api.search.v1.Area",
                "id": "44a3db8b-fa35-bb3d-bbdd-9a392690928b",
                "name": "Quận 6",
                "fullAddress": "Quận 6, TP. Hồ Chí Minh",
                "level": 3,
                "code": "",
                "parentId": "421741fc-2746-7e55-a5a4-eec23c50c986",
                "tags": "Quan 6, TP. Ho Chi Minh",
                "formattedAddress": ""
            },
            {
                "@type": "type.googleapis.com/vn.futa.vato.buslines.api.search.v1.Area",
                "id": "40b1b08a-04eb-1e90-a825-dc2cc925f46e",
                "name": "Quận 7",
                "fullAddress": "Quận 7, TP. Hồ Chí Minh",
                "level": 3,
                "code": "",
                "parentId": "421741fc-2746-7e55-a5a4-eec23c50c986",
                "tags": "Quan 7, TP. Ho Chi Minh",
                "formattedAddress": ""
            }
        ]
    }
}
```

Get all routes by filter

```
GET https://api-online.futabus.vn/vato/v1/search/routes?destAreaId=4d81d9de-5e7e-5623-a2c6-2c3306f59a23&destOfficeId=&originAreaId=421741fc-2746-7e55-a5a4-eec23c50c986&originOfficeId=&isReturn=false&isReturnTripLoad=false&fromDate=2026%2D03%2D30T00%3A00%3A00%2E000%2B07%3A00
```

Trong đó destAreaId,originAreaId  chính là districtId lấy từ API ở trên

Example response
```
{"requestId":"019d3f35-a074-7bf0-9b67-f1c0ae5d9770", "status":200, "error":null, "data":{"@type":"type.googleapis.com/vn.futa.vato.buslines.api.ListData", "items":[{"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.RouteSearchData", "routeId":"2021", "from":"", "to":""}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.RouteSearchData", "routeId":"1424", "from":"", "to":""}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.RouteSearchData", "routeId":"2020", "from":"", "to":""}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.RouteSearchData", "routeId":"1096", "from":"", "to":""}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.RouteSearchData", "routeId":"1420", "from":"", "to":""}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.RouteSearchData", "routeId":"2346", "from":"", "to":""}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.RouteSearchData", "routeId":"2473", "from":"", "to":""}]}}
```

Get trips by routes

```
POST https://api-online.futabus.vn/vato/v1/search/trip-by-route

example request body
{"minNumSeat":1,"channel":"mobile_app","fromDate":"2026-03-31T00:00:00.000+07:00","toDate":"2026-03-31T23:59:59.000+07:00","routeIds":["2021","1424","2020","1096","1420","2346","2473"],"sort":{"byPrice":"asc","byDepartureTime":"asc"},"page":0,"size":200}
```

Example response
```
{"requestId":"019d3f35-a295-7c40-866e-cc0089d6f5e8", "status":200, "error":null, "data":{"@type":"type.googleapis.com/vn.futa.vato.buslines.api.PaginationData", "page":0, "size":200, "total":15, "items":[{"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.TripSearchByRouteData", "tripId":"7283396", "departureTime":"2026-03-31T01:30:00Z", "rawDepartureTime":"08:30", "rawDepartureDate":"31-03-2026", "arrivalTime":"2026-03-31T18:30:00Z", "duration":1020, "seatTypeName":"Limousine", "price":410000, "emptySeatQuantity":12, "routeId":"1424", "distance":860, "wayId":"774", "maxSeatsPerBooking":5, "wayName":"BXAS - Ngã Tư Ga - Khu Công Nghệ Cao - Full Cao Tốc ( TP.HCM => Nha Trang ) - QL1A - Quảng Ngãi .", "wayNote":"Quý Khách đang chọn tuyến xe có lộ trình đi Cao Tốc từ TP. HCM đến Nha Trang, xe không nhận đón/ trả dọc đường quốc lộ 1A. Cần hỗ trợ thêm thông tin vui lòng liên hệ hotline 19006067.", "shuttleOption":null, "route":{"originCode":"TPHCM", "destCode":"QUANGNGAI", "originName":"BX An Sương", "destName":"Quảng Ngãi", "name":"An Suong - Quang Ngai", "originHubName":"Bến xe An Sương", "destHubName":"Bến Xe Quãng Ngãi"}, "seatTypeCode":"glm", "seatDiagramRefreshMs":15000}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.TripSearchByRouteData", "tripId":"7283852", "departureTime":"2026-03-31T03:00:00Z", "rawDepartureTime":"10:00", "rawDepartureDate":"31-03-2026", "arrivalTime":"2026-03-31T20:00:00Z", "duration":1020, "seatTypeName":"Limousine", "price":410000, "emptySeatQuantity":11, "routeId":"1424", "distance":860, "wayId":"774", "maxSeatsPerBooking":5, "wayName":"BXAS - Ngã Tư Ga - Khu Công Nghệ Cao - Full Cao Tốc ( TP.HCM => Nha Trang ) - QL1A - Quảng Ngãi .", "wayNote":"Quý Khách đang chọn tuyến xe có lộ trình đi Cao Tốc từ TP. HCM đến Nha Trang, xe không nhận đón/ trả dọc đường quốc lộ 1A. Cần hỗ trợ thêm thông tin vui lòng liên hệ hotline 19006067.", "shuttleOption":null, "route":{"originCode":"TPHCM", "destCode":"QUANGNGAI", "originName":"BX An Sương", "destName":"Quảng Ngãi", "name":"An Suong - Quang Ngai", "originHubName":"Bến xe An Sương", "destHubName":"Bến Xe Quãng Ngãi"}, "seatTypeCode":"glm", "seatDiagramRefreshMs":15000}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.TripSearchByRouteData", "tripId":"6423452", "departureTime":"2026-03-31T06:30:00Z", "rawDepartureTime":"13:30", "rawDepartureDate":"31-03-2026", "arrivalTime":"2026-03-31T23:30:00Z", "duration":1020, "seatTypeName":"Limousine", "price":410000, "emptySeatQuantity":14, "routeId":"2021", "distance":850, "wayId":"723", "maxSeatsPerBooking":5, "wayName":"BXMT - Ngã Tư Ga - Khu Công Nghệ Cao - Full Cao Tốc ( TP.HCM => Nha Trang ) - QL1A - Quảng Ngãi ", "wayNote":"Quý Khách đang chọn tuyến xe có lộ trình đi Cao Tốc từ TP. HCM đến Nha Trang, xe không nhận đón/ trả dọc đường quốc lộ 1A. Cần hỗ trợ thêm thông tin vui lòng liên hệ hotline 19006067.", "shuttleOption":null, "route":{"originCode":"TPHCM", "destCode":"QUANGNGAI", "originName":"TP.Hồ Chí Minh", "destName":"Quảng Ngãi", "name":"Mien Tay - Quang Ngai", "originHubName":"Bến Xe Miền Tây", "destHubName":"Bến Xe Quãng Ngãi"}, "seatTypeCode":"glm", "seatDiagramRefreshMs":15000}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.TripSearchByRouteData", "tripId":"7268728", "departureTime":"2026-03-31T08:30:00Z", "rawDepartureTime":"15:30", "rawDepartureDate":"31-03-2026", "arrivalTime":"2026-04-01T01:30:00Z", "duration":1020, "seatTypeName":"Limousine", "price":410000, "emptySeatQuantity":26, "routeId":"1096", "distance":840, "wayId":"167", "maxSeatsPerBooking":5, "wayName":"BXMĐ Mới => VP Suối Linh = > Cao Tốc Dầu Giây- Phan Thiết => QL1A => BX Quảng Ngãi", "wayNote":"Quý Khách đang chọn tuyến xe có lộ trình đi Cao Tốc  Dầu Giây - Phan Thiết. Lưu ý không nhận đón/ trả dọc đường quốc lộ 1A. Cần hỗ trợ thêm thông tin vui lòng liên hệ hotline 19006067.", "shuttleOption":null, "route":{"originCode":"TPHCM", "destCode":"QUANGNGAI", "originName":"TP.Hồ Chí Minh", "destName":"Quãng Ngãi", "name":"Mien Dong Moi - Quang Ngai", "originHubName":"BX Miền Đông Mới", "destHubName":"Bến Xe Quãng Ngãi"}, "seatTypeCode":"glm", "seatDiagramRefreshMs":15000}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.TripSearchByRouteData", "tripId":"7283853", "departureTime":"2026-03-31T09:00:00Z", "rawDepartureTime":"16:00", "rawDepartureDate":"31-03-2026", "arrivalTime":"2026-04-01T02:00:00Z", "duration":1020, "seatTypeName":"Limousine", "price":410000, "emptySeatQuantity":12, "routeId":"1424", "distance":860, "wayId":"774", "maxSeatsPerBooking":5, "wayName":"BXAS - Ngã Tư Ga - Khu Công Nghệ Cao - Full Cao Tốc ( TP.HCM => Nha Trang ) - QL1A - Quảng Ngãi .", "wayNote":"Quý Khách đang chọn tuyến xe có lộ trình đi Cao Tốc từ TP. HCM đến Nha Trang, xe không nhận đón/ trả dọc đường quốc lộ 1A. Cần hỗ trợ thêm thông tin vui lòng liên hệ hotline 19006067.", "shuttleOption":null, "route":{"originCode":"TPHCM", "destCode":"QUANGNGAI", "originName":"BX An Sương", "destName":"Quảng Ngãi", "name":"An Suong - Quang Ngai", "originHubName":"Bến xe An Sương", "destHubName":"Bến Xe Quãng Ngãi"}, "seatTypeCode":"glm", "seatDiagramRefreshMs":15000}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.TripSearchByRouteData", "tripId":"6423028", "departureTime":"2026-03-31T11:00:00Z", "rawDepartureTime":"18:00", "rawDepartureDate":"31-03-2026", "arrivalTime":"2026-04-01T03:00:00Z", "duration":960, "seatTypeName":"Limousine", "price":410000, "emptySeatQuantity":11, "routeId":"2473", "distance":870, "wayId":"770", "maxSeatsPerBooking":5, "wayName":"Bến Xe Ngã 4 Ga - Xa Lộ Hà Nội - Full Cao Tốc ( TP.HCM => Nha Trang ) - QL1A -Quy Nhơn- QL1A - Bến Xe Quảng Ngãi .", "wayNote":"", "shuttleOption":null, "route":{"originCode":"TPHCM", "destCode":"QUANGNGAI", "originName":"TP.Hồ Chí Minh", "destName":"Quảng Ngãi", "name":"Nga Tu Ga - Quang Ngai", "originHubName":"Bến xe Ngã 4 Ga", "destHubName":"Bến Xe Quãng Ngãi"}, "seatTypeCode":"glm", "seatDiagramRefreshMs":15000}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.TripSearchByRouteData", "tripId":"6423118", "departureTime":"2026-03-31T11:00:00Z", "rawDepartureTime":"18:00", "rawDepartureDate":"31-03-2026", "arrivalTime":"2026-04-01T04:00:00Z", "duration":1020, "seatTypeName":"Limousine", "price":410000, "emptySeatQuantity":19, "routeId":"1424", "distance":860, "wayId":"774", "maxSeatsPerBooking":5, "wayName":"BXAS - Ngã Tư Ga - Khu Công Nghệ Cao - Full Cao Tốc ( TP.HCM => Nha Trang ) - QL1A - Quảng Ngãi .", "wayNote":"Quý Khách đang chọn tuyến xe có lộ trình đi Cao Tốc từ TP. HCM đến Nha Trang, xe không nhận đón/ trả dọc đường quốc lộ 1A. Cần hỗ trợ thêm thông tin vui lòng liên hệ hotline 19006067.", "shuttleOption":null, "route":{"originCode":"TPHCM", "destCode":"QUANGNGAI", "originName":"BX An Sương", "destName":"Quảng Ngãi", "name":"An Suong - Quang Ngai", "originHubName":"Bến xe An Sương", "destHubName":"Bến Xe Quãng Ngãi"}, "seatTypeCode":"glm", "seatDiagramRefreshMs":15000}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.TripSearchByRouteData", "tripId":"6423695", "departureTime":"2026-03-31T11:00:00Z", "rawDepartureTime":"18:00", "rawDepartureDate":"31-03-2026", "arrivalTime":"2026-04-01T04:00:00Z", "duration":1020, "seatTypeName":"Limousine", "price":410000, "emptySeatQuantity":18, "routeId":"2021", "distance":850, "wayId":"723", "maxSeatsPerBooking":5, "wayName":"BXMT - Ngã Tư Ga - Khu Công Nghệ Cao - Full Cao Tốc ( TP.HCM => Nha Trang ) - QL1A - Quảng Ngãi ", "wayNote":"Quý Khách đang chọn tuyến xe có lộ trình đi Cao Tốc từ TP. HCM đến Nha Trang, xe không nhận đón/ trả dọc đường quốc lộ 1A. Cần hỗ trợ thêm thông tin vui lòng liên hệ hotline 19006067.", "shuttleOption":null, "route":{"originCode":"TPHCM", "destCode":"QUANGNGAI", "originName":"TP.Hồ Chí Minh", "destName":"Quảng Ngãi", "name":"Mien Tay - Quang Ngai", "originHubName":"Bến Xe Miền Tây", "destHubName":"Bến Xe Quãng Ngãi"}, "seatTypeCode":"glm", "seatDiagramRefreshMs":15000}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.TripSearchByRouteData", "tripId":"7272366", "departureTime":"2026-03-31T11:50:00Z", "rawDepartureTime":"18:50", "rawDepartureDate":"31-03-2026", "arrivalTime":"2026-04-01T04:50:00Z", "duration":1020, "seatTypeName":"Limousine", "price":410000, "emptySeatQuantity":12, "routeId":"1096", "distance":840, "wayId":"66", "maxSeatsPerBooking":5, "wayName":"BXMĐ Mới - Full Cao Tốc ( TP.HCM => Nha Trang )- QL1A - Quảng Ngãi.", "wayNote":"Quý Khách đang chọn tuyến xe có lộ trình đi Cao Tốc từ TP. HCM đến Nha Trang, xe không nhận đón/ trả dọc đường quốc lộ 1A. Cần hỗ trợ thêm thông tin vui lòng liên hệ hotline 19006067.", "shuttleOption":null, "route":{"originCode":"TPHCM", "destCode":"QUANGNGAI", "originName":"TP.Hồ Chí Minh", "destName":"Quãng Ngãi", "name":"Mien Dong Moi - Quang Ngai", "originHubName":"BX Miền Đông Mới", "destHubName":"Bến Xe Quãng Ngãi"}, "seatTypeCode":"glm", "seatDiagramRefreshMs":15000}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.TripSearchByRouteData", "tripId":"6422815", "departureTime":"2026-03-31T12:30:00Z", "rawDepartureTime":"19:30", "rawDepartureDate":"31-03-2026", "arrivalTime":"2026-04-01T05:30:00Z", "duration":1020, "seatTypeName":"Limousine", "price":410000, "emptySeatQuantity":18, "routeId":"1096", "distance":840, "wayId":"167", "maxSeatsPerBooking":5, "wayName":"BXMĐ Mới => VP Suối Linh = > Cao Tốc Dầu Giây- Phan Thiết => QL1A => BX Quảng Ngãi", "wayNote":"Quý Khách đang chọn tuyến xe có lộ trình đi Cao Tốc  Dầu Giây - Phan Thiết. Lưu ý không nhận đón/ trả dọc đường quốc lộ 1A. Cần hỗ trợ thêm thông tin vui lòng liên hệ hotline 19006067.", "shuttleOption":null, "route":{"originCode":"TPHCM", "destCode":"QUANGNGAI", "originName":"TP.Hồ Chí Minh", "destName":"Quãng Ngãi", "name":"Mien Dong Moi - Quang Ngai", "originHubName":"BX Miền Đông Mới", "destHubName":"Bến Xe Quãng Ngãi"}, "seatTypeCode":"glm", "seatDiagramRefreshMs":15000}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.TripSearchByRouteData", "tripId":"6538088", "departureTime":"2026-03-31T03:30:00Z", "rawDepartureTime":"10:30", "rawDepartureDate":"31-03-2026", "arrivalTime":"2026-03-31T23:30:00Z", "duration":1200, "seatTypeName":"Limousine", "price":510000, "emptySeatQuantity":4, "routeId":"2020", "distance":980, "wayId":"385", "maxSeatsPerBooking":5, "wayName":"BXMT - VP Ngã 4 Ga   - Cao Tốc LTDG - Cao Tốc DGPT - VP Phan Thiết - QL1A - BX Đà Nẵng", "wayNote":"Quý Khách đang chọn tuyến xe có lộ trình đi Cao Tốc- Long Thành - Dầu Giây - Phan Thiết. Lưu ý không nhận đón/ trả dọc đường quốc lộ 1A. Cần hỗ trợ thêm thông tin vui lòng liên hệ hotline 19006067.", "shuttleOption":null, "route":{"originCode":"TPHCM", "destCode":"DANANG", "originName":"TP.Hồ Chí Minh", "destName":"Đà Nẵng", "name":"Mien Tay - Da Nang", "originHubName":"Bến Xe Miền Tây", "destHubName":"Bến Xe Trung Tâm Đà Nẵng"}, "seatTypeCode":"glm", "seatDiagramRefreshMs":15000}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.TripSearchByRouteData", "tripId":"6537097", "departureTime":"2026-03-31T07:30:00Z", "rawDepartureTime":"14:30", "rawDepartureDate":"31-03-2026", "arrivalTime":"2026-04-01T03:30:00Z", "duration":1200, "seatTypeName":"Limousine", "price":510000, "emptySeatQuantity":21, "routeId":"1420", "distance":990, "wayId":"236", "maxSeatsPerBooking":5, "wayName":"BX An Sương - VP Ngã 4 Ga - VP Suối Linh - QL1A - BX Đà Nẵng", "wayNote":"", "shuttleOption":null, "route":{"originCode":"TPHCM", "destCode":"DANANG", "originName":"BX An Sương", "destName":"Đà Nẵng", "name":"An Suong - Da Nang", "originHubName":"Bến xe An Sương", "destHubName":"Bến Xe Trung Tâm Đà Nẵng"}, "seatTypeCode":"glm", "seatDiagramRefreshMs":15000}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.TripSearchByRouteData", "tripId":"6538057", "departureTime":"2026-03-31T08:00:00Z", "rawDepartureTime":"15:00", "rawDepartureDate":"31-03-2026", "arrivalTime":"2026-04-01T04:00:00Z", "duration":1200, "seatTypeName":"Limousine", "price":510000, "emptySeatQuantity":26, "routeId":"2020", "distance":980, "wayId":"385", "maxSeatsPerBooking":5, "wayName":"BXMT - VP Ngã 4 Ga   - Cao Tốc LTDG - Cao Tốc DGPT - VP Phan Thiết - QL1A - BX Đà Nẵng", "wayNote":"Quý Khách đang chọn tuyến xe có lộ trình đi Cao Tốc- Long Thành - Dầu Giây - Phan Thiết. Lưu ý không nhận đón/ trả dọc đường quốc lộ 1A. Cần hỗ trợ thêm thông tin vui lòng liên hệ hotline 19006067.", "shuttleOption":null, "route":{"originCode":"TPHCM", "destCode":"DANANG", "originName":"TP.Hồ Chí Minh", "destName":"Đà Nẵng", "name":"Mien Tay - Da Nang", "originHubName":"Bến Xe Miền Tây", "destHubName":"Bến Xe Trung Tâm Đà Nẵng"}, "seatTypeCode":"glm", "seatDiagramRefreshMs":15000}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.TripSearchByRouteData", "tripId":"6537066", "departureTime":"2026-03-31T11:15:00Z", "rawDepartureTime":"18:15", "rawDepartureDate":"31-03-2026", "arrivalTime":"2026-04-01T07:15:00Z", "duration":1200, "seatTypeName":"Limousine", "price":510000, "emptySeatQuantity":19, "routeId":"1420", "distance":990, "wayId":"236", "maxSeatsPerBooking":5, "wayName":"BX An Sương - VP Ngã 4 Ga - VP Suối Linh - QL1A - BX Đà Nẵng", "wayNote":"", "shuttleOption":null, "route":{"originCode":"TPHCM", "destCode":"DANANG", "originName":"BX An Sương", "destName":"Đà Nẵng", "name":"An Suong - Da Nang", "originHubName":"Bến xe An Sương", "destHubName":"Bến Xe Trung Tâm Đà Nẵng"}, "seatTypeCode":"glm", "seatDiagramRefreshMs":15000}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.TripSearchByRouteData", "tripId":"7268231", "departureTime":"2026-03-31T11:30:00Z", "rawDepartureTime":"18:30", "rawDepartureDate":"31-03-2026", "arrivalTime":"2026-04-01T07:30:00Z", "duration":1200, "seatTypeName":"Limousine", "price":510000, "emptySeatQuantity":24, "routeId":"2020", "distance":980, "wayId":"385", "maxSeatsPerBooking":5, "wayName":"BXMT - VP Ngã 4 Ga   - Cao Tốc LTDG - Cao Tốc DGPT - VP Phan Thiết - QL1A - BX Đà Nẵng", "wayNote":"Quý Khách đang chọn tuyến xe có lộ trình đi Cao Tốc- Long Thành - Dầu Giây - Phan Thiết. Lưu ý không nhận đón/ trả dọc đường quốc lộ 1A. Cần hỗ trợ thêm thông tin vui lòng liên hệ hotline 19006067.", "shuttleOption":null, "route":{"originCode":"TPHCM", "destCode":"DANANG", "originName":"TP.Hồ Chí Minh", "destName":"Đà Nẵng", "name":"Mien Tay - Da Nang", "originHubName":"Bến Xe Miền Tây", "destHubName":"Bến Xe Trung Tâm Đà Nẵng"}, "seatTypeCode":"glm", "seatDiagramRefreshMs":15000}], "others":[]}}
```

Get seats by tripId

```
GET /vato/v1/search/seat-diagram/7283396
Host: api-online.futabus.vn
```
Response
```
HTTP/2 200 OK
Server: openresty
Date: Mon, 30 Mar 2026 14:46:14 GMT
Content-Type: application/json

{"requestId":"019d3f35-a85a-722a-aa15-512ec8824931", "status":200, "error":null, "data":{"@type":"type.googleapis.com/vn.futa.vato.buslines.api.PaginationData", "page":0, "size":33, "total":33, "items":[{"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394445", "name":"B17", "status":[], "columnNo":5, "rowNo":13, "floor":"up", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394446", "name":"A01", "status":[1], "columnNo":1, "rowNo":1, "floor":"down", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394447", "name":"A02", "status":[1], "columnNo":5, "rowNo":1, "floor":"down", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394448", "name":"A03", "status":[1], "columnNo":1, "rowNo":2, "floor":"down", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394449", "name":"A06", "status":[1], "columnNo":1, "rowNo":3, "floor":"down", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394450", "name":"A09", "status":[1], "columnNo":1, "rowNo":4, "floor":"down", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394451", "name":"A12", "status":[1], "columnNo":1, "rowNo":5, "floor":"down", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394452", "name":"A15", "status":[1], "columnNo":1, "rowNo":6, "floor":"down", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394453", "name":"A16", "status":[1], "columnNo":3, "rowNo":6, "floor":"down", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394454", "name":"A17", "status":[1], "columnNo":5, "rowNo":6, "floor":"down", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394455", "name":"A11", "status":[1], "columnNo":5, "rowNo":4, "floor":"down", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394456", "name":"A08", "status":[1], "columnNo":5, "rowNo":3, "floor":"down", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394457", "name":"A05", "status":[1], "columnNo":5, "rowNo":2, "floor":"down", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394458", "name":"A04", "status":[1], "columnNo":3, "rowNo":2, "floor":"down", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394459", "name":"A07", "status":[1], "columnNo":3, "rowNo":3, "floor":"down", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394460", "name":"A10", "status":[1], "columnNo":3, "rowNo":4, "floor":"down", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394461", "name":"A13", "status":[1], "columnNo":3, "rowNo":5, "floor":"down", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394462", "name":"B01", "status":[1], "columnNo":1, "rowNo":8, "floor":"up", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394463", "name":"B02", "status":[1], "columnNo":5, "rowNo":8, "floor":"up", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394464", "name":"B03", "status":[1], "columnNo":1, "rowNo":9, "floor":"up", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394465", "name":"B04", "status":[1], "columnNo":3, "rowNo":9, "floor":"up", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394466", "name":"B05", "status":[1], "columnNo":5, "rowNo":9, "floor":"up", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394468", "name":"B07", "status":[], "columnNo":3, "rowNo":10, "floor":"up", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394469", "name":"B08", "status":[], "columnNo":5, "rowNo":10, "floor":"up", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394470", "name":"B09", "status":[], "columnNo":1, "rowNo":11, "floor":"up", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394471", "name":"B10", "status":[], "columnNo":3, "rowNo":11, "floor":"up", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394472", "name":"B11", "status":[], "columnNo":5, "rowNo":11, "floor":"up", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394473", "name":"B12", "status":[], "columnNo":1, "rowNo":12, "floor":"up", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394474", "name":"B13", "status":[], "columnNo":3, "rowNo":12, "floor":"up", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394475", "name":"B14", "status":[], "columnNo":5, "rowNo":12, "floor":"up", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394476", "name":"B15", "status":[], "columnNo":1, "rowNo":13, "floor":"up", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279394477", "name":"B16", "status":[], "columnNo":3, "rowNo":13, "floor":"up", "price":410000, "property":0}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.SeatDiagramData", "seatId":"279413305", "name":"B06", "status":[], "columnNo":1, "rowNo":10, "floor":"up", "price":410000, "property":0}], "others":[]}}
```

GET department-in-way by routeId

```
GET /vato/v1/search/department-in-way/774?routeId=1424 HTTP/2
Host: api-online.futabus.vn
```

Response
```
HTTP/2 200 OK
Server: openresty
Date: Mon, 30 Mar 2026 14:46:19 GMT
Content-Type: application/json

{"requestId":"019d3f35-bbb9-7d05-8154-aba604238ed7", "status":200, "error":null, "data":{"@type":"type.googleapis.com/vn.futa.vato.buslines.api.PaginationData", "page":0, "size":8, "total":8, "items":[{"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.DepartmentInWayData", "departmentId":"441", "departmentName":"NGUYỄN ẢNH THỦ", "departmentAddress":"186A Nguyễn Ảnh Thủ, Phường Hiệp Thành, Quận 12, TP.Hồ Chí Minh", "wardId":"", "wardName":"", "districtId":"", "districtName":"", "provinceId":"", "provinceName":"", "timeAtDepartment":-45, "passing":true, "isShuttleService":false, "note":"", "latitude":0, "longitude":0, "pointKind":2, "presentBeforeMinutes":-45, "openingTime":"", "closingTime":"", "onlineOpeningTime":"", "onlineClosingTime":"", "isActive":true}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.DepartmentInWayData", "departmentId":"518", "departmentName":"VĨNH LỘC - AN SƯƠNG", "departmentAddress":"390 Đường số 01, KDC Vĩnh Lộc, P.Bình Tân, TPHCM", "wardId":"", "wardName":"", "districtId":"", "districtName":"", "provinceId":"", "provinceName":"", "timeAtDepartment":-45, "passing":true, "isShuttleService":false, "note":"", "latitude":0, "longitude":0, "pointKind":2, "presentBeforeMinutes":-45, "openingTime":"", "closingTime":"", "onlineOpeningTime":"", "onlineClosingTime":"", "isActive":true}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.DepartmentInWayData", "departmentId":"537", "departmentName":"749 QUANG TRUNG (HCM)", "departmentAddress":"749 QUANG TRUNG (HCM)", "wardId":"", "wardName":"", "districtId":"", "districtName":"", "provinceId":"", "provinceName":"", "timeAtDepartment":-45, "passing":true, "isShuttleService":false, "note":"", "latitude":0, "longitude":0, "pointKind":2, "presentBeforeMinutes":-45, "openingTime":"", "closingTime":"", "onlineOpeningTime":"", "onlineClosingTime":"", "isActive":true}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.DepartmentInWayData", "departmentId":"423", "departmentName":"Xa Lộ Đại Hàn", "departmentAddress":"2389 Quốc Lộ 1A, phường Tân Hưng Thuận, Quận 12, TP. Hồ Chí Minh", "wardId":"", "wardName":"", "districtId":"", "districtName":"", "provinceId":"", "provinceName":"", "timeAtDepartment":-30, "passing":true, "isShuttleService":false, "note":"", "latitude":0, "longitude":0, "pointKind":2, "presentBeforeMinutes":-30, "openingTime":"", "closingTime":"", "onlineOpeningTime":"", "onlineClosingTime":"", "isActive":true}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.DepartmentInWayData", "departmentId":"231", "departmentName":"BX An Sương", "departmentAddress":"Bến Xe An Sương, Quốc Lộ 22, Ấp Đông Lân, Bà Điểm, Hóc Môn, TP Hồ Chí Minh", "wardId":"", "wardName":"", "districtId":"", "districtName":"", "provinceId":"", "provinceName":"", "timeAtDepartment":0, "passing":true, "isShuttleService":true, "note":"", "latitude":10.84406, "longitude":106.6138, "pointKind":0, "presentBeforeMinutes":-15, "openingTime":"", "closingTime":"", "onlineOpeningTime":"", "onlineClosingTime":"", "isActive":true}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.DepartmentInWayData", "departmentId":"408", "departmentName":"NGÃ TƯ GA", "departmentAddress":"BX NGÃ 4 GA", "wardId":"", "wardName":"", "districtId":"", "districtName":"", "provinceId":"", "provinceName":"", "timeAtDepartment":15, "passing":true, "isShuttleService":false, "note":"", "latitude":0, "longitude":0, "pointKind":2, "presentBeforeMinutes":0, "openingTime":"", "closingTime":"", "onlineOpeningTime":"", "onlineClosingTime":"", "isActive":true}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.DepartmentInWayData", "departmentId":"222", "departmentName":"BX Nam Tuy Hòa", "departmentAddress":"VP Bến Xe Nam Tuy Hòa: 507 Nguyễn Văn Linh, phường Phú Lâm, TP.Tuy Hòa, Phú Yên", "wardId":"", "wardName":"", "districtId":"", "districtName":"", "provinceId":"", "provinceName":"", "timeAtDepartment":600, "passing":true, "isShuttleService":false, "note":"", "latitude":13.04065, "longitude":109.3116, "pointKind":2, "presentBeforeMinutes":585, "openingTime":"", "closingTime":"", "onlineOpeningTime":"", "onlineClosingTime":"", "isActive":true}, {"@type":"type.googleapis.com/vn.futa.vato.buslines.api.search.v1.DepartmentInWayData", "departmentId":"83", "departmentName":"Quảng Ngãi", "departmentAddress":"Số 02 Trần Khánh Dư, phường Nghĩa Chánh, Thành phố Quảng Ngãi", "wardId":"", "wardName":"", "districtId":"", "districtName":"", "provinceId":"", "provinceName":"", "timeAtDepartment":930, "passing":true, "isShuttleService":false, "note":"", "latitude":15.11112, "longitude":108.8219, "pointKind":1, "presentBeforeMinutes":915, "openingTime":"", "closingTime":"", "onlineOpeningTime":"", "onlineClosingTime":"", "isActive":true}], "others":[]}}
```

Booking

```
POST /vato/v1/booking/reservation HTTP/2
Host: api-online.futabus.vn
```

```
{"passenger":{"custName":"Đỗ Quốc Tuấn","loginMobile":"0367717714","custEmail":"tuandoquoc28@gmail.com","custSn":"","custMobile":"0367717714"},"ticketInfo":[{"seats":[{"seatId":"279394470"}],"dropoff":{"lng":108.8219,"address":"Số 02 Trần Khánh Dư, phường Nghĩa Chánh, Thành phố Quảng Ngãi","officeId":"83","lat":15.11112,"name":"Quảng Ngãi","timeAtDepartment":915,"type":3},"tripId":"7283396","pickup":{"officeId":"231","name":"BX An Sương","timeAtDepartment":-15,"address":"Bến Xe An Sương, Quốc Lộ 22, Ấp Đông Lân, Bà Điểm, Hóc Môn, TP Hồ Chí Minh","type":3,"lat":10.844060000000001,"lng":106.6138}}],"channel":"mobile_app"}
```

Response
```
HTTP/2 200 OK
Server: openresty
Date: Mon, 30 Mar 2026 14:46:28 GMT
Content-Type: application/json
Content-Length: 372

{"requestId":"019d3f35-dcf5-73e5-9343-13bbcc5e8037", "status":200, "error":null, "data":{"@type":"type.googleapis.com/vn.futa.vato.buslines.api.SingleData", "item":{"@type":"type.googleapis.com/vn.futa.vato.buslines.api.booking.v1.BookingTicketResponse", "id":"019d3f35-deed-7d1c-8e4d-533e6e066e35", "code":"ETINRS", "totalPrice":410000, "discount":0, "isHoliday":false}}}
```

Get ticket detail by id

```
GET /vato/v1/user/ticket/detail?id=019d3f35-deed-7d1c-8e4d-533e6e066e35&code=ETINRS&phone=0367717714 HTTP/2
Host: api-online.futabus.vn
```

Response
```
HTTP/2 200 OK
Server: openresty
Date: Mon, 30 Mar 2026 14:46:29 GMT
Content-Type: application/json

{"requestId":"019d3f35-e356-7725-abc3-01a6391f8dff", "status":200, "error":null, "data":{"@type":"type.googleapis.com/vn.futa.vato.buslines.api.SingleData", "item":{"@type":"type.googleapis.com/vn.futa.vato.buslines.api.user.v1.Ticket", "id":"019d3f35-deed-7d1c-8e4d-533e6e066e35", "createdAt":"2026-03-30T14:46:28.333852Z", "updatedAt":"2026-03-30T14:46:28.333852Z", "code":"ETINRS", "status":"booking", "wayType":"one_way", "bookingChannel":"mobile_app", "paymentMethod":0, "promotionCode":"", "originAmount":410000, "discountAmount":0, "refundAmount":0, "referenceUserId":"2170167", "referencePhone":"0367717714", "remarks":"", "ticketSeats":[{"id":"019d3f35-def2-778e-ae1b-757f3108ba4b", "routeId":"1424", "seatType":"outbound", "routeName":"An Suong - Quang Ngai", "departureTime":"2026-03-31T01:30:00Z", "seatId":"279394470", "seatName":"B09", "status":"initial", "originAmount":410000, "discountAmount":0, "refundAmount":0, "paidAt":"0001-01-01T00:00:00Z", "exportedAt":"0001-01-01T00:00:00Z", "cancelledAt":"0001-01-01T00:00:00Z", "refundedAt":"0001-01-01T00:00:00Z", "kind":"", "originCode":"", "originName":"An Suong ", "destCode":"", "destName":" Quang Ngai", "distanceKm":0, "durationMs":"1020", "wayId":"774", "wayName":"", "pickupType":"office", "pickupId":"", "pickupName":"BX An Sương", "pickupInfo":{"address":"Bến Xe An Sương, Quốc Lộ 22, Ấp Đông Lân, Bà Điểm, Hóc Môn, TP Hồ Chí Minh", "appointment_time":0, "hub_id":0, "lat":10.84406, "lng":106.6138, "name":"BX An Sương", "office_id":"231", "pickup_id":0, "place_id":"", "time_at_department":-15, "type":3, "zone_id":0}, "dropOffType":"office", "dropOffId":"", "dropOffName":"", "dropOffInfo":{"address":"Số 02 Trần Khánh Dư, phường Nghĩa Chánh, Thành phố Quảng Ngãi", "appointment_time":0, "hub_id":0, "lat":15.11112, "lng":108.8219, "name":"Quảng Ngãi", "office_id":"83", "pickup_id":0, "place_id":"", "time_at_department":915, "type":3, "zone_id":0}, "extraData":{"cancelFeeRatio":30, "checkinCode":"279394470", "invoiceCode":"793944703", "timeAcceptCancel":"2026-03-30T01:30:00Z", "timeIgnoreCancel":"2026-03-30T01:50:00Z"}, "tripId":"7283396"}], "ticketPassengers":[{"id":"019d3f35-def3-7dc6-b2cc-33fddf531133", "createdAt":"2026-03-30T14:46:28.339907Z", "updatedAt":"2026-03-30T14:46:28.339907Z", "ticketId":"019d3f35-deed-7d1c-8e4d-533e6e066e35", "ticketSeatId":"019d3f35-def2-778e-ae1b-757f3108ba4b", "name":"Đỗ Quốc Tuấn", "phone":"0367717714", "email":"tuandoquoc28@gmail.com", "identityCard":"", "properties":{}}], "ticketExtras":[{"id":"019d3f35-def4-7ea9-aa7d-21d3a3f28ce7", "createdAt":"2026-03-30T14:46:28.340955Z", "updatedAt":"2026-03-30T14:46:29.913167Z", "ticketId":"019d3f35-deed-7d1c-8e4d-533e6e066e35", "data":{"cancelFeeRatio":30, "isHoliday":false, "timeAcceptCancel":"2026-03-30T01:30:00Z", "timeExpiredPayment":"2026-03-30T21:51:28+07:00", "timeIgnoreCancel":"2026-03-30T01:50:00Z", "vehiclePlate":""}}]}}}
```