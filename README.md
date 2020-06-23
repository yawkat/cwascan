cwascan
=======

Scan for exposure notification BLE messages and insert them into postgres. Adapted from [this gist](https://gist.github.com/thetzel/398c5c504a4616732e78c991e2478e52). Needs this patch to paypal gatt (also in the gist):

```diff
diff --git a/adv.go b/adv.go
index 787ff33..d4f0cc5 100644
--- a/adv.go
+++ b/adv.go
@@ -127,7 +127,12 @@ func (a *Advertisement) unmarshall(b []byte) error {
                case typeManufacturerData:
                        a.ManufacturerData = make([]byte, len(d))
                        copy(a.ManufacturerData, d)
-               // case typeServiceData16,
+               case typeServiceData16:
+                       // https://gist.github.com/thetzel/398c5c504a4616732e78c991e2478e52
+                       var s ServiceData
+                       s.UUID = UUID{d[:2]}
+                       s.Data = d[2:]
+                       a.ServiceData = append(a.ServiceData, s)
                // case typeServiceData32,
                // case typeServiceData128:
                default:
```

db schema:

```sql
CREATE TABLE rpik (
    "timestamp" timestamp without time zone DEFAULT timezone('utc'::text, now()),
    rpik bytea,
    metadata bytea,
    rssi integer
);
```