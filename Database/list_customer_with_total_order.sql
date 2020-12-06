select c.customer_name,
       c.email,
       c.phone_number,
       c.dob,
       c.sex,
       sum(case when o.customer_id is null then 0 else 1 end) total_orders
from customers c
         left join orders o on c.customer_id = o.customer_id
group by c.customer_id;