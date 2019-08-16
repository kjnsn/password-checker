# NIST Bad Password Checker

This is a simple service that checks a password against a dictionary of common
passwords. The service then returns a response informing the caller whether or
not the password is common.

## Implementation

The server loads a dictionary of passwords, then creates a bloom filter in order
to efficiently check the input password for presence in the dictionary. False
positives are possible but unlikely.

## Source Code Headers

Every file containing source code must include copyright and license
information. This includes any JS/CSS files that you might be serving out to
browsers. (This is to help well-intentioned people avoid accidental copying that
doesn't comply with the license.)

Apache header:

    Copyright 2019 Google LLC

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

        https://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.

## Disclaimer

This is not an officially supported Google product.
