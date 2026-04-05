-- By default the system sets this to 100. Due to our pooling/sinkers/etc... we can hit this limit in integration tests.
-- Long term this can go away once we clean up some of the sinkers/tests, but for
ALTER SYSTEM SET max_connections = 500;