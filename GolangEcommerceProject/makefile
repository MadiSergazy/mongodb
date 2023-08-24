SIGNUP_API := http://localhost:8000/users/signup
LOGIN_API := http://localhost:8000/users/login
ADD_PRODUCT_API := http://localhost:8000/admin/addproduct
VIEW_PRODUCTS_API := http://localhost:8000/users/productview
SEARCH_PRODUCTS_API := http://localhost:8000/users/search
ADD_TO_CART_API := http://localhost:8000/addtocart
REMOVE_FROM_CART_API := http://localhost:8000/removefromcart
LIST_CART_API := http://localhost:8000/listcart
ADD_ADDRESS_API := http://localhost:8000/addaddress
EDIT_HOME_ADDRESS_API := http://localhost:8000/edithomeaddress
EDIT_WORK_ADDRESS_API := http://localhost:8000/editworkaddress
DELETE_ADDRESSES_API := http://localhost:8000/deleteaddresses
PLACE_ORDER_API := http://localhost:8000/placeorder
INSTANT_BUY_API := http://localhost:8000/instantbuy

.PHONY: signup login add-product view-products search-products add-to-cart remove-from-cart list-cart add-address edit-home-address edit-work-address delete-addresses place-order instant-buy

signup:
	curl -X POST -H "Content-Type: application/json" -d "{\"first_name\": \"Akhil\", \"last_name\": \"Sharma\", \"email\": \"akhil@gmail.com\", \"password\": \"akhilsharma\", \"phone\": \"+4534545435\"}" $(SIGNUP_API)

login:
	curl -X POST -H "Content-Type: application/json" -d "{\"email\": \"akhil@gmail.com\", \"password\": \"akhilsharma\"}" $(LOGIN_API)

add-product:
	curl -X POST -H "Content-Type: application/json" -d '{
		"product_name": "Alienware x15",
		"price": 2500,
		"rating": 10,
		"image": "alienware.jpg"
	}' $(ADD_PRODUCT_API)

view-products:
	curl $(VIEW_PRODUCTS_API)

search-products:
	curl "$(SEARCH_PRODUCTS_API)?name=al"

add-to-cart:
	curl "$(ADD_TO_CART_API)?id=xxxproduct_idxxx&userID=xxxxxxuser_idxxxxxx"

remove-from-cart:
	curl "$(REMOVE_FROM_CART_API)?id=xxxxxxx&userID=xxxxxxxxxxxx"

list-cart:
	curl "$(LIST_CART_API)?id=xxxxxxuser_idxxxxxxxxxx"

add-address:
	curl -X POST -H "Content-Type: application/json" -d '{
		"house_name": "white house",
		"street_name": "white street",
		"city_name": "washington",
		"pin_code": "332423432"
	}' "$(ADD_ADDRESS_API)?id=user_id**\*\***\***\*\***"

edit-home-address:
	curl -X PUT "$(EDIT_HOME_ADDRESS_API)?id=xxxxxxxxxxuser_idxxxxxxxxxxxxxxx"

edit-work-address:
	curl -X PUT "$(EDIT_WORK_ADDRESS_API)?id=xxxxxxxxxxuser_idxxxxxxxxxxxxxxx"

delete-addresses:
	curl "$(DELETE_ADDRESSES_API)?id=xxxxxxxxxuser_idxxxxxxxxxxxxx"

place-order:
	curl "$(PLACE_ORDER_API)?id=xxuser_idxxx"

instant-buy:
	curl "$(INSTANT_BUY_API)?userid=xxuser_idxxx&pid=xxxxproduct_idxxxx"
