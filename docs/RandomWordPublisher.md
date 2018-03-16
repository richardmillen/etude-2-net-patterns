# Random Word Publisher Example

## Publisher

this application publishes random(ish) words in a few different languages to subscribers.
the server works on a fire & forget basis so late joiners will never receive previously
sent messages unless the 'word-lvc' app is used.
subscribers are able to receive all updates for one or more of the languages (topics).

## Subscriber

this application subscribes to a never ending series of live messages
(random words) sent by a publisher (either 'word-pub' or 'word-lvc').
the app can subscribe to words in one or more languages (topics);
'eng', 'fra' & 'esp'.

if this app subscribes directly to a 'word-pub' instance then it will receive
messages as & when they are sent by the publisher. on the other hand, if it
subscribes to a 'word-lvc' instance then it will receive an update (the most
recent relevant message/s) immediately after a connection is established instead
of having to wait for the next 'live' message.

## Last-Value Cache

this application sits between a publisher and one or more subscribers. it subscribes
to the publisher, caching the last value sent for each topic as it arrives. it then
immediately sends relevant cached value(s) to new subscribers as they connect.




