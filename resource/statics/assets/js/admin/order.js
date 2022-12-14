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
    if (request.orderStatus != undefined) {
        params.push({"name": "orderStatus", "value": request.orderStatus});
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

        var html = "<thead><tr><th width='250px'>ID/外部单号</th><th width='350px'>金额信息</th><th width='350px'>订单商品</th><th width='300px'>地址</th><th width='350px'>其他信息</th><th width='100px'>操作</th></tr></thead><tbody id='tbody'></tbody>"
        var tr = "<tr><td>{id}</td><td>{amountInfo}</td><td>{orderProducts}</td><td>{addressInfo}</td><td>{extra}</td><td>{operationBtn}</td></tr>"
        list = jsonResult.data.list
        $.each(list, function(index, item) {
            html += tr.replace(/{(\w+)}/g, function(_, $1) {
                if ($1 == 'id') {
                    id = "" 
                    if (item.orderStatus == 0) {
                        id += "<span class='mr-2'> <span class='badge-dot badge-primary'>";
                    } else if (item.orderStatus == 10) {
                        id += "<span class='mr-2'> <span class='badge-dot badge-brand'>";
                    } else if (item.orderStatus == 15) {
                        id += "<span class='mr-2'><span class='badge-dot  badge-secondary'>";
                    } else if (item.orderStatus == 20) {
                        id += "<span class='mr-2'><span class='badge-dot  badge-secondary'>";
                    } else if (item.orderStatus == 30) {
                        id += "<span class='mr-2'><span class='badge-dot badge-success'>"
                    } else if (item.orderStatus == 40) {
                        id += "<span class='mr-2'><span class='badge-dot badge-dark'>"
                    } else {
                        id += "<span class='mr-2'>  <span class='badge-dot badge-light'>"
                    }
                    return id += "</span>" + item.orderStatusText + "<code>" + item['id'] + "</code><br><code>" + item['outOrderNo'] + "</code><br>" + item["createdTime"]

                } else if ($1 == 'amountInfo') {
                    var amount = "";
                    amount += "总金额：" + item.totalAmount / 100 + "(货币:" + item.currency + ")<br>"
                    amount += "折扣金额：" + item.totalDiscountAmount / 100 + "<br>"
                    amount += "实付金额：<code>" + item.realTotalAmount / 100 + "(运费:" + item.postalFee / 100 + ")</code><br>"
                    return amount
                } else if ($1 == 'orderProducts') {
                    var products = "";
                    if (item["productItems"]) {
                        for (k = 0; k < item["productItems"].length; k++) {
                            productItem = item['productItems'][k]
                            if (productItem && productItem.product && productItem.sku) {
                                if (k != 0) {
                                    products += "<hr>"
                                }
                                products += "<b>" + productItem.product.title + "-" + productItem.sku.title + " * " + "<code>" + productItem.count + "</code></b>  <img src='" + productItem.thumbnail + "' height='55px' alt='加载...' />"
                            }
                        }
                    } else if (item["thumbnail"]) {
                        for (k = 0; k < item["thumbnail"].length; k++) {
                            if (k != 0) {
                                products += "&&<img src='" + item["thumbnail"][k] + "' height='55px' alt='加载失败' />"
                            } else {
                                products += "<img src='" + item["thumbnail"][k] + "' height='55px' alt='加载失败' />"
                            }
                        }
                    }
                    return products
                } else if ($1 == 'addressInfo') {
                    return item[$1]
                } else if ($1 == 'extra') {
                    var extra = ""
                    if (item["extra"]) {
                        for (k = 0; k < item["extra"].length; k++) {
                            if (k != 0) {
                                extra += "<hr>"
                            }
                            extra += "*" + item["extra"][k]["name"] + ": " + item["extra"][k]["item"]
                        }
                    }
                    return extra
                    // if (!!item.extra && item.extra.length > 0) {
                    //     return JSON.stringify(item.extra)
                    // }
                } else if ($1 == 'operationBtn') {
                    operationBtn = "<input type='button' class='btn btn-outline-success dropdown-toggle' data-toggle='dropdown' value='操作'>"
                    operationBtn += "<ul class='dropdown-menu'>"
                    operationBtn += "<li ><a class='btn btn-outline-brand' href='javascript:changOrderStatus(" + item.id + "," + 10 + ")'>设为已下单</a></li>"
                    operationBtn += "<li ><a class='btn btn-outline-brand' href='javascript:changOrderStatus(" + item.id + "," + 15 + ")'>设为部分到货</a></li>"
                    operationBtn += "<li ><a class='btn btn-outline-brand' href='javascript:changOrderStatus(" + item.id + "," + 20 + ")'>设为已到货</a></li>"
                    operationBtn += "<li ><a class='btn btn-outline-brand' href='javascript:changOrderStatus(" + item.id + "," + 30 + ")'>设为已发货</a></li>"
                    operationBtn += "<li ><a class='btn btn-outline-brand' href='javascript:changOrderStatus(" + item.id + "," + 40 + ")'>设为已取消</a></li>"
                    operationBtn += "<li ><a class='btn btn-outline-brand' href='javascript:editExtraForm(" + JSON.stringify(item).replace(/\”/g,"'") + ")'>编辑</a></li>"
                    operationBtn += "</ul>"
                    return operationBtn

                }
                return item[$1]

            });

        });
        $(".order-list").append(html)
        // $(".dropdown-toggle").dropdown('toggle');
        listPaginate(jsonResult.data.total)
        return
    }).catch(function(error) {
        console.log(error);
    });
}
function orderStatistic() {

    var url = "/admin/order/statistics";

    options = {
        method: 'get',
        headers: {
            "Content-Type": "application/json charset=utf-8",
        },
    }

    commonFetch(url, options).then(function (jsonResult) {
        if (jsonResult.code != 0) {
            alert("获取订单数据失败, err: " + jsonResult.msg)
            return
        }

        var statistic = "<div><label>"
        statistic += "<div style='float: left; border: 1px solid blue; padding: 5px; margin=2px'><font color='blue'>待下单：{waitPay}</font></div>"
        statistic += "<div style='float: left; border: 1px solid blue; padding: 5px;  margin=2px'><font color='blue'> 待到货：{waitGot}</font></div>"
        statistic += "<div style='float: left; border: 1px solid blue; padding: 5px; margin=2px'><font color='blue'> 待发货：{waitDeliver}</font></div>"
        // statistic += "<div style='float: left; border: 1px solid blue; padding: 5px; margin=2px'><font color='blue'> 已经发货：{delivered}</font></div>"
        statistic += "</label></div>"

        item = jsonResult.data.item

        html = statistic.replace(/{(\w+)}/g, function(_, $1) {
            if ($1 == "waitPay" || $1 == "waitGot" || $1 == "waitDeliver" || $1 == "delivered") {
                return item[$1]
            } 
            return ""
        });
        $(".order-statistic").append(html)
    })
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

function editExtraForm(order)
{
    idHtml = "<input type='hidden' id='exitExtraOrderId' name='orderId' value='" + order.id + "'>"

    $("#orderAddressExtra").html(order.addressInfo);
    $("#orderExtend").html(JSON.stringify(order.extra));

    $("#edit-extra").modal("show")
    $("#form-hidden-items").append(idHtml)
}

function editExtra()
{
    url = "/admin/order/edit/extra";
    body = {}
    if ($("#orderExtend").val()) {
        body.orderExtra = JSON.parse($("#orderExtend").val());
    } 
    if ($("#orderAddressExtra").val()) {
        body.addressExtra = $("#orderAddressExtra").val();
    }
    if (!body.orderExtra && !body.addressExtra) {
        alert("请不要提交空信息")
        return 
    }
    body.orderId = parseInt($("#exitExtraOrderId").val());
    options = {
        method: 'post',
        headers: {
            "Content-Type": "application/json charset=utf-8",
        },
        body: JSON.stringify(body)
    }
    console.log(options.body)
    commonFetch(url, options).then(function (jsonResult) {
        if (jsonResult.code != 0) {
            alert("修改订单扩展失败, err: " + jsonResult.msg)
            return
        }
        location.reload()
    })

}

function changOrderStatus(orderId, orderStatus)
{
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