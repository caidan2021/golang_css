/*
 * @Date: 2022-08-31 11:08:56
 */

function orderList() {

    var url = "/admin/order/list?";

    var params = [];
    request = GetRequest()
    if (request.orderNo != undefined) {
        params.push({"name": "orderNo", "value": request.orderNo});
    }
    if (request.page != undefined) {
        params.push({"name": "page", "value": request.page});
    } 
    url = ReplaceOrAddRequest(params, url)

    options = {
        method: 'get',
        headers: {
            "Content-Type": "application/json charset=utf-8",
        },
    }

    commonFetch(url, options).then(function (jsonResult) {
        if (jsonResult.code != 0) {
            alert("获取订单失败, err: " + jsonResult.msg)
            return
        }

        var html = "<thead><tr><th width='5%'>ID</th><th>封面</th><th width='10%'>订单No</th><th width='10%'>状态</th><th width='10%'>创建时间</th><th width='20%'>操作</th></tr></thead><tbody id='tbody'></tbody>"
        var tr = "<tr><td>{id}</td><td>{thumbnail}</td><td>{outOrderNo}</td><td>{orderStatusText}</td><td>{createdTime}</td><td>{operationBtn}</td></tr>"

        list = jsonResult.data.list
        $.each(list, function(index, item) {
            html += tr.replace(/{(\w+)}/g, function(_, $1) {
                if ($1 == 'thumbnail') {
                    var thumbnails = "";
                    for (k = 0; k < item[$1].length; k++) {
                        if (k != 0) {
                            thumbnails += "&&<img src='" + item[$1][k] + "' height='75px' alt='加载失败' />"
                        } else {
                            thumbnails += "<img src='" + item[$1][k] + "' height='75px' alt='加载失败' />"
                        }
                    }
                    return thumbnails
                }
                if ($1 == 'orderStatusText') {
                    if (item.orderStatus == 0) {
                        orderStatusText = "<span class='mr-2'> <span class='badge-dot badge-primary'>";
                    } else if (item.orderStatus == 10) {
                        orderStatusText = "<span class='mr-2'> <span class='badge-dot badge-brand'>";
                    } else if (item.orderStatus == 20) {
                        orderStatusText = "<span class='mr-2'><span class='badge-dot  badge-secondary'>";
                    } else if (item.orderStatus == 30) {
                        orderStatusText = "<span class='mr-2'><span class='badge-dot badge-success'>"
                    } else {
                        orderStatusText = "<span class='mr-2'>  <span class='badge-dot badge-light'>"
                    }
                    return orderStatusText += "</span>" + item.orderStatusText + "</span>"
                }

                if ($1 == 'operationBtn') {
                    operationBtn = "<input type='button' class='btn btn-outline-success dropdown-toggle' data-toggle='dropdown' value='操作'>"
                    operationBtn += "<ul class='dropdown-menu'>"
                    operationBtn += "<li ><a class='btn btn-outline-secondary' href='javascript:changOrderStatus(" + item.id + "," + 10 + ")'>设为已下单</a></li>"
                    operationBtn += "<li ><a class='btn btn-outline-brand' href='javascript:changOrderStatus(" + item.id + "," + 20 + ")'>设为到货</a></li>"
                    operationBtn += "<li ><a class='btn btn-outline-success' href='javascript:changOrderStatus(" + item.id + "," + 30 + ")'>设为已发货</a></li>"
                    operationBtn += "</ul>"
                    return operationBtn
                }
                return item[$1]
            });

        });
        $(".order-list").append(html)
        $(".dropdown-toggle").dropdown('toggle');
        listPaginate(jsonResult.data.total)
    }).catch(function(error) {
        console.log(error);
    });
}

function createOrder()
{
    imgs = new Array();
    $(".hidden_img").each(function () {
        imgs.push($(this).val())
    });
    data = {
        "thumbnails": imgs,
        "thirdPartyFlag": $("#thirdPartyFlag").val(),
        'outOrderNo': $("#outOrderNo").val(),
    }
    url = "/admin/order/create"
    options = {
        method: 'post',
        headers: {
            "Content-Type": "application/json charset=utf-8",
        },
        body: JSON.stringify(data)
    }
    commonFetch(url, options).then(function (jsonResult) {
        if (jsonResult.code != 0) {
            alert("创建订单失败, err: " + jsonResult.msg)
            return
        }
        window.location = "/admin/view/admin_order_list"
    })

}

function changOrderStatus(orderId, orderStatus)
{
    console.log(orderId,orderStatus)
    url = "/admin/order/change/status";
    options = {
        method: 'post',
        headers: {
            "Content-Type": "application/json charset=utf-8",
        },
        body: JSON.stringify({
            "id": orderId,
            "orderStatus": orderStatus
        })
    }
    commonFetch(url, options).then(function (jsonResult) {
        if (jsonResult.code != 0) {
            alert("修改订单状态失败, err: " + jsonResult.msg)
            return
        }
        location.reload()
    })
}

function createOrderForm()
{
    $("#createOrderForm").modal("show")
}