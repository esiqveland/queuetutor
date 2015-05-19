# NSQ aka 
 - Not Stable Queue, 
 - Naming Sucks Queue, 
 - New Simple Queue


### Guarantees:
    - messages are not durable (by default)
        - can overflow messages to disk (hint: set overflow to 0!)
    - _at least once_ delivery
    - messages received are _un-ordered_
    - "consumers eventually find all topic producers"
        - consequence/feature in avoiding single point of failure in a distributed fashion


## A small example to what you can use queueing and nsq for!

1. a producer of a 'task', here a user submits an application for some product or feature from the frontend
2. consumers of this 'task' or more 
    - a we want a confirmation email to be generated and sent from the emailing system
    - we want to record this signup as a metric
    - store it in the db

producer:
    run  

