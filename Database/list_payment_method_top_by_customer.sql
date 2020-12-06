select pm.method_name,
       pm.code,
       sum(IF(o.payment_method_id is null, 0, 1)) most_used
from payment_methods pm
    left join orders o on pm.payment_method_id = o.payment_method_id
    left join customers c on o.customer_id = c.customer_id
group by pm.payment_method_id
order by most_used desc;