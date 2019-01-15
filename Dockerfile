FROM centos:7

RUN yum -y install epel-release git-core

RUN yum -y install maven unzip java-1.8.0-openjdk
