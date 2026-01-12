# Solution 

#### Locking Strategy
In a concurrent environment, multiple users may attempt to claim the same coupon simultaneously,
which is introduced a race condition issue. </br>

To prevent this, I serialize access to the claim process using a distributed lock using Redis.
When a user attempts to claim a coupon, the system tries to set a unique key `claim:coupon:{coupon_name}` in Redis;
if this key is successfully set, the process "wins" the lock and proceeds the claim process by inserting an entry in
user_claims table and decrement the `coupon.remaining_amount `by 1. Concurrent requests will detect that the key already
exists, forcing them to pause for `500ms` and retry up to `3` times. </br>

This approach eliminates the race conditions in the coupon claim process by enforcing that only 1 claim is processed,
while maintaining the high performance of the API. Please check the documentation directory for the detailed solution. </br>

### Consideration 
- Suitable for very high throughput with in-memory access using Redis, preventing database bottlenecks or deadlocks.
- Non-blocking process with retry mechnism: 3 retries and 500ms retry delay. </br>
  The maximum latency of the API would be only 1.5s for the worst-case scenario, allowing the system to fail fast 
  and provide immediate response rather than hanging indefinitely.
- Safety & Deadlock Prevention: Reliability is guaranteed through a Time-to-Live (TTL) mechanism; 
  In case of instance crashes or failures, the lock automatically expires after 20 seconds, 
  therefore, the claim process eventually becomes available again.

### Stress Test
I made a stress test to test the efficiency of my solution in stress_test/main.go by sending 100 requests concurrently 
to claim a coupon with 5 quotas. 

The result of the benchmark is 5 success and 95 failed requests with average response time ~500ms. 

### Alternative
Message Queue
</br>**Pros** 
- The API returns an immediate response and the heavy process operates in the background. 
  The message broker ensures requests are processed sequentially by consumers, protecting the database from traffic spikes.
-  Built-in retry mechanisms ensure transient failures are handled automatically. 
   In case failure still persists after retrials, failed messages are safely routed to a Dead Letter Queue (DLQ), 
   ensuring zero data loss and allowing for manual inspection/reprocessing.

**Cons :** 
- Requires maintaining additional infrastructure, such as producer/consumer application and broker setup.  
  It needs to handle idempotency issue to ensuring a message only processes once.
- Informing succes outcome to the users would be more difficult, either through notification, Webscoket, or another service

### Further Development 
1. Create reclaim functionality to cancel the claimed coupons in case of failed or cancelled transaction in the real-world scenario.
2. Create eligible functionality to determine whether a user is able to use the coupon before claiming it. 

### Author
Chris Christian 
