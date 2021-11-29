FROM    ubuntu

# Install python3
RUN     apt-get update
RUN     apt-get install -y python3

# Copy html
COPY . .

CMD ["python3", "-m", "http.server", "8000"]