# Backend Engineer programming exercise

Thanks for trying our programming exercise!

With this document you should have received three executable files:

+ `ascii_linux_x64`: executable binary compiled for linux x64 platforms
+ `ascii_darwin_x64`: executable binary compiled for MacOS x64 platforms
+ `ascii_windows_x64`: executable binary compiled for Windows x64 platforms

Please contact us immediately if:

+ the executable for your platform is missing
+ you think there are problems with the provided executable files

## The Challenge

The proposed challenge is to implement an ASCII art hosting webservice. The solution  should consist in a REST HTTP API.

The provided executable will simulate clients interactions with the webservice. They will connect through port `4444` on `localhost`.  Once connected, clients will be sending three kinds of requests:

1. **Image registration**: The client requests an image registration for further upload. When registering an image, the client provides its SHA256 hash for further reference. Registering a preexisting image should result in an error (`409 Conflict`).
2. **Image chunks upload**: The client splits the image content in a sequence of chunks and uploads them. It sends each chunk separately as a JSON payload. Each chunk has an ID indicating its position in the sequence.
3. **Downloading the complete image**: The client downloads the image back from the webservice. It then computes the downloaded image hash and compares it to the registered image. It is expected that an image could be downloaded multiple times.

A single image related sequence of events will always follow the one mentioned above.
The executable's output for a single image upload sequence will look like the following:

```bash
./ascii_linux_x64 -files 1
time="2019-02-28T11:55:15+01:00" level=info msg="registering image with server" image_sha256=8a99030199b315fe8e4cf93d93478facdf1801a0ddb0d9bc1325961597a42a3f
time="2019-02-28T11:55:15+01:00" level=info msg="uploading image chunks" chunks_count=3 image_sha256=8a99030199b315fe8e4cf93d93478facdf1801a0ddb0d9bc1325961597a42a3f
time="2019-02-28T11:55:15+01:00" level=info msg="chunk upload: OK" chunk_id=1 chunk_size=256 image_sha256=8a99030199b315fe8e4cf93d93478facdf1801a0ddb0d9bc1325961597a42a3f
time="2019-02-28T11:55:16+01:00" level=info msg="chunk upload: OK" chunk_id=0 chunk_size=256 image_sha256=8a99030199b315fe8e4cf93d93478facdf1801a0ddb0d9bc1325961597a42a3f
time="2019-02-28T11:55:16+01:00" level=info msg="chunk upload: OK" chunk_id=2 chunk_size=187 image_sha256=8a99030199b315fe8e4cf93d93478facdf1801a0ddb0d9bc1325961597a42a3f
time="2019-02-28T11:55:16+01:00" level=fatal msg="succesfully retrieved image" image_sha256=8a99030199b315fe8e4cf93d93478facdf1801a0ddb0d9bc1325961597a42a3f
```

## The Protocol

Our executable expects your HTTP API to implement the following endpoints:

+ **Registering an image**:
  + **method**: `POST`
  + **URI**: `/image`
  + **Content-Type**: `application/json`
  + **Request Body**:

      ```json
      {
      "sha256": "abc123easyasdoremi...",
      "size": 123456,
      "chunk_size": 256
      }
      ```

  + **Responses**:
    | Code                       |              Description           |
    |----------------------------|------------------------------------|
    | 201 Created                | Image succesfully registered       |
    | 409 Conflict               | Image already exists               |
    | 400 Bad Request            | Malformed request                  |
    | 415 Unsupported Media Type | Unsupported payload format         |

+ **Uploading an image chunk**:
  + **method**: `POST`
  + **URI**: `/image/<sha256>/chunks`
  + **Content-Type**: `application/json`
  + **Request Body**:

      ```json
      {
      "id": 1,
      "size": 256,
      "data": "8   888   , 888    Y888 888 888    ,ee 888 888 888 888 ...",
      }
      ```

  + **Responses**:
    | Code          |              Description           |
    |---------------|------------------------------------|
    | 201 Created   | Chunk succesfully uploaded         |
    | 409 Conflict  | Chunk already exists               |
    | 404 Not Found | Image not found                    |

+ **Downloading an image**:
  + **method**: `GET`
  + **URI**: `/image/<sha256>`
  + **Accept**: `text/plain`
  + **Responses**:
    | Code          |              Description           |
    |---------------|------------------------------------|
    | 200 OK        | Image succesfully downloaded       |
    | 404 Not Found | Image not found                    |

  + **Nota Bene**: Notice this endpoint returns plaintext, not JSON. It also expects to download the whole file rather than separate chunks.

+ **Errors**:
  + **Accept**: `application/json`
  + **Response body**:

    ```json
    {
      "code": "HTTP Error code. e.g: 400 Bad Request",
      "message": "Error message. e.g: Malformed request; reason: chunk id field is missing."
    }
    ```

## The client executable

You can tune the executable's behavior using the following command-line options:

```bash
  -chunksize int
        size of chunks used (default 256)
  -files int
        Amount of files to generate and send to the host (default 300)
  -host string
        host to send the requests to (default "localhost")
  -port int
        port to use when sending requests to the host (default 4444)
  -log-format string
        set structured logging format; possible values: text, json (default "text")
  -log-level string
        set logging level; possible values: debug, info, warn, error (default "info")
  -seed int
        set the seed used to produce randomness; providing a value will allow reproducible runs (default -1)
```

As a matter of fact we will use these same options to test the behavior of your server. We recommend that you use these options to ensure your solution covers corner cases.

**Nota Bene**: The executable supports the `-seed` option to ease debugging. Providing it an integer value ensures reproducible client runs.

## Your Solution

The challenge is designed so that a candidate can implement the
solution in a limited amount of time, thus feel
free to use your preferred framework, libraries, and tools.

Your solution should provide a `Dockerfile`. We will use it to build and run your server. It should expose port `4444`. To test your solution we will run the following commands:

```bash
docker build -t recruitment/<candidate> .
docker run -d -p 4444:4444 recruitment/<candidate>
./ascii_<platform>_x64
```

**N.B**: We are open to solutions using `docker-compose` if its usage is justified. We would then run the following commands: ```docker-compose build && docker-compose up -d```.

We are open to review solutions in any non-exotic language. Yet, most team members are comfortable with: *Go*, *Rust*, *Python*, *Node.js* and *C/C++*. Although not proficient in those languages, we are also open to review solutions written in *Java*, *Scala*, *Ruby*, *Clojure* or *Haskell*. We recommend you use the language you are the most comfortable with.
**We expect you to ensure your solution works against the client executable before sending it to us**. The first thing we will do with your code is to run it against the provided client. You can, thus, consider it as a test suite in charge of providing you very early feedback.

To test your solution, first make sure you have your server
running and listening on `http://localhost:4444`. Also ensure our executable has execution rights (`chmod +x`), and simply run it:

```bash
./ascii_<platform>_x64
```

This will start our testing client, and will immediately start communicating with your server. You will know it passed when it outputs:

```text
SUCCESS!
```

## Assumptions

You can assume your code will be ran on a modern, powerful machine equipped with:

+ A multi-core CPU
+ 8GB+ of RAM
+ a SSD disk

## Assessment Criteria

We expect you to write code you would consider production-ready.
This means we expect your code to be well-factored, without needless
duplication and following good practices.

What we will look at:

+ If your code fulfils the requirements, and runs against the
supplied client.
+ How clean is your design and implementation, how easy it is to
understand and maintain your code.
+ How your server behaves under stress: CPU usage, Memory Usage, IO usage, system calls...
