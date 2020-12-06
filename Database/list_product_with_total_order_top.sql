select p.product_name,
       p.basic_price,
       sum(case when od.product_id is null then 0 else 1 end) number_orders
from products p
         left join order_details od on p.product_id = od.product_id
         left join orders o on od.order_id = o.order_id
group by p.product_id
order by number_orders desc;