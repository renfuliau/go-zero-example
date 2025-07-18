
§
order_service.protoorder"j
CreateOrderReq
userId (RuserId&
items (2.order.OrderItemRitems
address (	Raddress"[
	OrderItem
	productId (R	productId
quantity (Rquantity
price (Rprice"E
CreateOrderResp
orderId (	RorderId
success (Rsuccess"'
GetOrderReq
orderId (	RorderId"À
GetOrderResp
orderId (	RorderId
userId (RuserId&
items (2.order.OrderItemRitems
status (	Rstatus 
totalAmount (RtotalAmount
	createdAt (	R	createdAt2z
Order<
CreateOrder.order.CreateOrderReq.order.CreateOrderResp3
GetOrder.order.GetOrderReq.order.GetOrderRespB	Z./orderbproto3