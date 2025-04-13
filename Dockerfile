# Copyright 2023 LiveKit, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.24-alpine AS builder

ARG TARGETPLATFORM
ARG TARGETARCH
RUN echo building for "$TARGETPLATFORM"

WORKDIR /workspace

# Install git
RUN apk add --no-cache git

# Clone the LiveKit repo and checkout latest version
RUN git clone https://github.com/livekit/livekit.git .
RUN git fetch --all --tags
RUN git checkout tags/v1.10.1 -b v1.10.1


# Download Go modules
RUN go mod download

# Build the server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH GO111MODULE=on go build -a -o livekit-server ./cmd/server


FROM alpine

COPY --from=builder /workspace/livekit-server /livekit-server
# Copy the livekit.yaml config into the image
COPY livekit.yaml /etc/livekit.yaml

# Run the binary with the config
ENTRYPOINT ["/livekit-server", "--config", "/etc/livekit.yaml"]
