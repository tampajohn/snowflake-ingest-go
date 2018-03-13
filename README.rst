Snowflake Ingest Service Go SDK
*********************************

.. image:: http://img.shields.io/:license-Apache%202-brightgreen.svg
    :target: http://www.apache.org/licenses/LICENSE-2.0.txt

The Snowflake Ingest Service SDK allows users to ingest files
into their Snowflake data warehouse in a programmatic fashion via key-pair
authentication.

Prerequisites
=============

Go 1.8+
-------

The Snowflake Ingest Service SDK can only be used with Go 1.8 or higher. Previous versions
may work

A 2048-bit RSA key pair
-----------------------
Snowflake Authentication for the Ingest Service requires creating a 2048 bit
RSA key pair and, registering the public key with Snowflake. For detailed instructions,
please visit the relevant `Snowflake Documentation Page <docs.snowflake.net>`_.


Adding as a Dependency
======================
You can add the Snowflake Ingest Service SDK by running the following command

.. code-block:: bash

    go get github.com/tampajohn/