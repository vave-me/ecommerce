"use client";
import React from 'react';
import Messaging from "./Messaging";
import LazyNATSProvider from "../../../components/Utils/LazyNATSProvider";
export default function MessagesPage() {
    return <LazyNATSProvider> <Messaging/></LazyNATSProvider>
}