#!/bin/sh

systemctl daemon-reload
systemctl enable prometheus-arista-exporter
systemctl start prometheus-arista-exporter
