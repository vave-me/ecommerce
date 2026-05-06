"use client";
export const dynamic = 'force-dynamic';
import React, {useMemo, useState} from "react";
import * as Yup from "yup";
import {useTranslations} from "next-intl";
import {useForm} from "react-hook-form";
import {yupResolver} from "@hookform/resolvers/yup";
import {useRouter, useSearchParams} from "next/navigation";
import {useAuth} from "../../../context/AuthContext";
import Script from 'next/script';
import Link from 'next/link';
import Logo from '../../../components/Header/Logo';
import styles from "./RegisterForm.module.css";
import { suggestAddress } from '../../../api/geocodingApi';

export default function RegisterForm() {
    const t = useTranslations("LoginForm");
    const router = useRouter();
    const searchParams = useSearchParams();
    const {signUpWithCredentials, signInWithGoogle} = useAuth();
    
    // Get redirect URL from query parameters, default to home page
    const redirectUrl = searchParams.get('redirect') || searchParams.get('returnTo') || '/';
    
    /* UI state */
    const [error, setError] = useState("");
    const [addressQuery, setAddressQuery] = useState("");
    const [addressSuggestions, setAddressSuggestions] = useState([]);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [googleLoaded, setGoogleLoaded] = useState(false);
    const googleButtonRef = React.useRef(null);
    
    /* Yup schema for registration only */
    const schema = useMemo(() => Yup.object({
        email: Yup.string()
            .email(t("errors.invalidEmail"))
            .required(t("errors.required")),
        password: Yup.string()
            .required(t("errors.required"))
            .min(6, t("errors.passwordTooShort", "Password must be at least 6 characters")),
        confirmPassword: Yup.string()
            .oneOf([Yup.ref("password")], t("errors.passwordsMustMatch"))
            .required(t("errors.required")),
        firstName: Yup.string().required(t("errors.required")),
        lastName: Yup.string().required(t("errors.required")),
        userName: Yup.string().required(t("errors.required")),
        address: Yup.string().required(t("errors.required"))
    }), [t]);
    
    /* React-Hook-Form */
    const {
        register,
        handleSubmit,
        setValue,
        formState: {errors}
    } = useForm({
        resolver: yupResolver(schema),
        defaultValues: {
            email: "", 
            password: "", 
            confirmPassword: "",
            firstName: "", 
            lastName: "", 
            userName: "", 
            address: ""
        }
    });
    
    /* Submit handler - register only */
    const onSubmit = async data => {
        setError("");
        setIsSubmitting(true);
        try {
            await signUpWithCredentials({
                email: data.email,
                password: data.password,
                firstName: data.firstName,
                lastName: data.lastName,
                userName: data.userName,
                address: data.address,
                role: 'customer'  // Always set role to customer for regular sign-ups
            }, redirectUrl);
        } catch (err) {
            setError(err.response?.data?.message || t("errors.authFailed"));
        } finally {
            setIsSubmitting(false);
        }
    };
    
    /* Address suggestions */
    const handleAddressChange = async e => {
        const val = e.target.value;
        setValue("address", val);
        setAddressQuery(val);
        if (val.length > 2) {
            try {
                const result = await suggestAddress(val);
                if (result.suggestionAddress && result.suggestionAddress.length > 0) {
                    const addresses = result.suggestionAddress.map(item => item.suggestedAddress);
                    setAddressSuggestions(addresses);
                } else {
                    setAddressSuggestions([]);
                }
            } catch (err) {
                // Error: 'Error fetching address suggestions:', err...
                setAddressSuggestions([]);
            }
        } else {
            setAddressSuggestions([]);
        }
    };
    
    const handleGoogleCredentialResponse = async (response) => {
        try {
            if (!process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID) {
                throw new Error("Google Client ID is not configured. Please check your environment variables.");
            }
            const idToken = response.credential;
            if (!idToken) {
                throw new Error("No ID token received from Google");
            }
            await signInWithGoogle(idToken, redirectUrl);
        } catch (err) {
            setError(err.message || t("errors.authFailed"));
        }
    };
    
    /* Handle navigation back */
    const handleGoBack = () => {
        if (redirectUrl && redirectUrl !== '/') {
            router.push(redirectUrl);
        } else {
            router.push('/');
        }
    };

    /* Render */
    return (
        <div className={styles.container}>
            {/* Navigation Header */}
            <nav className={styles.navHeader}>
                <button 
                    onClick={handleGoBack}
                    className={styles.closeButton}
                    aria-label="Go back"
                    type="button"
                >
                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                        <line x1="18" y1="6" x2="6" y2="18"></line>
                        <line x1="6" y1="6" x2="18" y2="18"></line>
                    </svg>
                </button>
                <div className={styles.brandLogo}>
                    <Logo size="default" />
                </div>
            </nav>
            
            {/* Google Sign-In Script */}
            {typeof window !== 'undefined' && (
                <Script 
                    src="https://accounts.google.com/gsi/client" 
                    strategy="afterInteractive" 
                    onLoad={() => {
                        if (!process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID || !window.google) {
                            return;
                        }
                        
                        setTimeout(() => {
                            try {
                                google.accounts.id.initialize({
                                    client_id: process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID,
                                    callback: handleGoogleCredentialResponse,
                                    auto_select: false,
                                    cancel_on_tap_outside: true,
                                    context: "signup"
                                });
                                
                                if (googleButtonRef.current) {
                                    google.accounts.id.renderButton(
                                        googleButtonRef.current,
                                        { 
                                            theme: "outline", 
                                            size: "large", 
                                            type: "standard",
                                            width: "100%",
                                            text: "signup_with",
                                            shape: "rectangular",
                                            logo_alignment: "center"
                                        }
                                    );
                                    setGoogleLoaded(true);
                                }
                            } catch (error) {
                                // Error: 'Google Sign-In error:', error...
                                setError("Failed to initialize Google Sign-In. Please try again later.");
                            }
                        }, 100);
                    }}
                />
            )}
            
            <section className={styles.section}>
                <div className={styles.hero}>
                    <h1>{t("createAccount")}</h1>
                    <p>Join thousands of users who are already discovering and sharing on our platform. Create your account in just a few steps.</p>
                </div>
                
                <div className={styles.formWrapper}>
                    <form onSubmit={handleSubmit(onSubmit)} className={styles.form}>
                        {error && <div className={styles.errorMessage}>{error}</div>}
                        
                        {/* firstName */}
                        <div className={styles.inputContainer}>
                            <label htmlFor="firstName" className={styles.inputLabel}>
                                {t("firstName")}
                            </label>
                            <div className={styles.inputWrapper}>
                                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={styles.inputIcon}>
                                    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                                    <circle cx="12" cy="7" r="4"></circle>
                                </svg>
                                <input
                                    id="firstName"
                                    placeholder={t("placeholders.firstName")}
                                    className={`${styles.input} ${styles.withIcon}`}
                                    {...register("firstName")}
                                    aria-invalid={!!errors.firstName}
                                />
                            </div>
                            {errors.firstName && (
                                <div className="mt-1 text-sm text-red-500">
                                    {errors.firstName.message}
                                </div>
                            )}
                        </div>
                        
                        {/* lastName */}
                        <div className={styles.inputContainer}>
                            <label htmlFor="lastName" className={styles.inputLabel}>
                                {t("lastName")}
                            </label>
                            <div className={styles.inputWrapper}>
                                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={styles.inputIcon}>
                                    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                                    <circle cx="12" cy="7" r="4"></circle>
                                </svg>
                                <input
                                    id="lastName"
                                    placeholder={t("placeholders.lastName")}
                                    className={`${styles.input} ${styles.withIcon}`}
                                    {...register("lastName")}
                                    aria-invalid={!!errors.lastName}
                                />
                            </div>
                            {errors.lastName && (
                                <div className="mt-1 text-sm text-red-500">
                                    {errors.lastName.message}
                                </div>
                            )}
                        </div>
                        
                        {/* userName */}
                        <div className={styles.inputContainer}>
                            <label htmlFor="userName" className={styles.inputLabel}>
                                {t("userName")}
                            </label>
                            <div className={styles.inputWrapper}>
                                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={styles.inputIcon}>
                                    <path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"></path>
                                    <circle cx="12" cy="7" r="4"></circle>
                                </svg>
                                <input
                                    id="userName"
                                    placeholder={t("placeholders.userName")}
                                    className={`${styles.input} ${styles.withIcon}`}
                                    {...register("userName")}
                                    aria-invalid={!!errors.userName}
                                />
                            </div>
                            {errors.userName && (
                                <div className="mt-1 text-sm text-red-500">
                                    {errors.userName.message}
                                </div>
                            )}
                        </div>
                        
                        {/* address */}
                        <div className={styles.inputContainer}>
                            <label htmlFor="address" className={styles.inputLabel}>
                                {t("address")}
                            </label>
                            <div className={styles.inputWrapper}>
                                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={styles.inputIcon}>
                                    <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path>
                                    <circle cx="12" cy="10" r="3"></circle>
                                </svg>
                                <input
                                    id="address"
                                    placeholder={t("placeholders.address")}
                                    className={`${styles.input} ${styles.withIcon}`}
                                    value={addressQuery}
                                    onChange={handleAddressChange}
                                    aria-invalid={!!errors.address}
                                />
                            </div>
                            {errors.address && (
                                <div className="mt-1 text-sm text-red-500">
                                    {errors.address.message}
                                </div>
                            )}
                            {addressSuggestions.length > 0 && (
                                <ul className={styles.suggestionsList}>
                                    {addressSuggestions.map((s, i) => (
                                        <li
                                            key={i}
                                            className={styles.suggestionItem}
                                            onClick={() => {
                                                setValue("address", s);
                                                setAddressQuery(s);
                                                setAddressSuggestions([]);
                                            }}
                                        >
                                            {s}
                                        </li>
                                    ))}
                                </ul>
                            )}
                        </div>
                        
                        {/* email */}
                        <div className={styles.inputContainer}>
                            <label htmlFor="email" className={styles.inputLabel}>
                                {t("email")}
                            </label>
                            <div className={styles.inputWrapper}>
                                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={styles.inputIcon}>
                                    <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path>
                                    <polyline points="22,6 12,13 2,6"></polyline>
                                </svg>
                                <input
                                    id="email"
                                    type="email"
                                    placeholder={t("placeholders.email")}
                                    className={`${styles.input} ${styles.withIcon}`}
                                    {...register("email")}
                                    aria-invalid={!!errors.email}
                                />
                            </div>
                            {errors.email && (
                                <div className="mt-1 text-sm text-red-500">
                                    {errors.email.message}
                                </div>
                            )}
                        </div>
                        
                        {/* password */}
                        <div className={styles.inputContainer}>
                            <label htmlFor="password" className={styles.inputLabel}>
                                {t("password")}
                            </label>
                            <div className={styles.inputWrapper}>
                                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={styles.inputIcon}>
                                    <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                                    <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
                                </svg>
                                <input
                                    id="password"
                                    type="password"
                                    placeholder={t("placeholders.password")}
                                    className={`${styles.input} ${styles.withIcon}`}
                                    {...register("password")}
                                    aria-invalid={!!errors.password}
                                />
                            </div>
                            {errors.password && (
                                <div className="mt-1 text-sm text-red-500">
                                    {errors.password.message}
                                </div>
                            )}
                        </div>
                        
                        {/* confirmPassword */}
                        <div className={styles.inputContainer}>
                            <label htmlFor="confirmPassword" className={styles.inputLabel}>
                                {t("confirmPassword")}
                            </label>
                            <div className={styles.inputWrapper}>
                                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={styles.inputIcon}>
                                    <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                                    <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
                                </svg>
                                <input
                                    id="confirmPassword"
                                    type="password"
                                    placeholder={t("placeholders.confirmPassword")}
                                    className={`${styles.input} ${styles.withIcon}`}
                                    {...register("confirmPassword")}
                                    aria-invalid={!!errors.confirmPassword}
                                />
                            </div>
                            {errors.confirmPassword && (
                                <div className="mt-1 text-sm text-red-500">
                                    {errors.confirmPassword.message}
                                </div>
                            )}
                        </div>
                        
                        {/* submit */}
                        <button
                            type="submit"
                            disabled={isSubmitting}
                            className={styles.submitButton}
                        >
                            {isSubmitting && (
                                <svg className="animate-spin -ml-1 mr-2 h-5 w-5 text-white inline-block" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                                </svg>
                            )}
                            {t("createAccount")}
                        </button>
                        
                        <div className={styles.divider}>
                            <span>{t("or")}</span>
                        </div>
                        
                        <div className={styles.googleButtonWrapper}>
                            {!googleLoaded && (
                                <button type="button" className={styles.googleButtonPlaceholder}>
                                    <svg className={styles.googleIcon} viewBox="0 0 24 24">
                                        <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                                        <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                                        <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                                        <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
                                    </svg>
                                    <span>Sign up with Google</span>
                                </button>
                            )}
                            <div 
                                ref={googleButtonRef}
                                id="googleSignInButton" 
                                className={styles.googleButton} 
                                style={{ display: googleLoaded ? 'flex' : 'none' }}
                            />
                        </div>
                        
                        <div className={styles.switchMode}>
                            {t("alreadyHaveAccount")}{" "}
                            <Link href="/login">
                                {t("signIn")}
                            </Link>
                        </div>
                    </form>
                </div>
            </section>
        </div>
    );
}