"use client";
import React from "react";
import { useParams } from "next/navigation";
import PublicProfile from "../../../../../components/Profile/PublicProfile";
export default function ProfileIdPage() {
    const params = useParams();
    const { id } = params;
    return <PublicProfile userId={id} />;
} 