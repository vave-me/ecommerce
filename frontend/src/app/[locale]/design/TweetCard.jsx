import React, { useState } from 'react';
import { Heart, MessageCircle, Share2, Bookmark, Flag, MoreHorizontal, Check, Clock, MapPin, Camera } from '@/icons';
const SimplifiedPostCard = () => {
    // Sample post data
    const post = {
        author: {
            name: "Sarah Johnson",
            username: "@sarahjdesigns",
            avatar: "/api/placeholder/36/36",
            verified: true
        },
        content: "Just finished the redesign for our client's e-commerce platform! It's amazing how small UI improvements can lead to significant increases in conversion rates. #UXDesign #UI",
        image: "/api/placeholder/500/300",
        timeAgo: "2h ago",
        location: "Berlin",
        likes: 127,
        comments: 23,
        hashtags: ["UXDesign", "UI", "ConversionOptimization"]
    };
    // Interactive states
    const [liked, setLiked] = useState(false);
    const [saved, setSaved] = useState(false);
    const [likesCount, setLikesCount] = useState(post.likes);
    // Handle like action
    const handleLike = () => {
        if (liked) {
            setLikesCount(likesCount - 1);
        } else {
            setLikesCount(likesCount + 1);
        }
        setLiked(!liked);
    };
    return (
        <div className="bg-white rounded-lg shadow-sm border border-gray-100 overflow-hidden max-w-md mx-auto">
            {/* Post header */}
            <div className="flex items-center p-3">
                {/* Author avatar */}
                <div className="relative mr-2">
                    <img
                        src={post.author.avatar}
                        alt={post.author.name}
                        className="w-9 h-9 rounded-full object-cover"
                    />
                    {post.author.verified && (
                        <div className="absolute -bottom-0.5 -right-0.5 bg-blue-500 rounded-full p-0.5">
                            <Check size={8} className="text-white" />
                        </div>
                    )}
                </div>
                {/* Author info */}
                <div className="flex-1 min-w-0">
                    <div className="flex items-center">
                        <span className="font-medium text-gray-900 text-sm">{post.author.name}</span>
                        <span className="text-gray-500 text-xs ml-1">{post.author.username}</span>
                    </div>
                    <div className="flex items-center text-gray-500 text-xs">
                        <Clock size={12} className="mr-0.5" />
                        <span>{post.timeAgo}</span>
                        {post.location && (
                            <>
                                <span className="mx-1">•</span>
                                <MapPin size={12} className="mr-0.5" />
                                <span>{post.location}</span>
                            </>
                        )}
                    </div>
                </div>
                {/* Post options */}
                <button className="p-1 text-gray-500 hover:text-gray-700 rounded-full">
                    <MoreHorizontal size={16} />
                </button>
            </div>
            {/* Post content */}
            <div className="px-3 pb-2">
                <p className="text-sm text-gray-800">{post.content}</p>
            </div>
            {/* Post image with overlay buttons */}
            <div className="relative">
                <img
                    src={post.image}
                    alt="Post attachment"
                    className="w-full object-cover"
                />
                {/* Image interaction buttons */}
                <div className="absolute top-2 right-2 flex flex-col gap-2">
                    {/* Like button */}
                    <button
                        className="bg-black bg-opacity-60 hover:bg-opacity-70 p-2 rounded-full text-white transition-all active:scale-95"
                        onClick={handleLike}
                    >
                        <Heart
                            size={18}
                            className={liked ? "text-red-500 fill-current" : ""}
                        />
                    </button>
                    {/* Comment button */}
                    <button className="bg-black bg-opacity-60 hover:bg-opacity-70 p-2 rounded-full text-white transition-all active:scale-95">
                        <MessageCircle size={18} />
                    </button>
                </div>
            </div>
            {/* Engagement stats */}
            <div className="px-3 pt-2 flex items-center text-xs text-gray-500">
                <span className="mr-3">{likesCount} likes</span>
                <span>{post.comments} comments</span>
            </div>
            {/* Action buttons */}
            <div className="flex items-center justify-between px-3 py-2 border-t border-gray-100 mt-2">
                <div className="flex space-x-1">
                    {/* Save button */}
                    <button
                        className={`p-1.5 rounded-full hover:bg-gray-100 transition-colors ${saved ? 'text-indigo-600' : 'text-gray-500'}`}
                        onClick={() => setSaved(!saved)}
                    >
                        <Bookmark size={18} className={`${saved ? 'fill-current' : ''}`} />
                    </button>
                    {/* Share button */}
                    <button className="p-1.5 rounded-full hover:bg-gray-100 transition-colors text-gray-500">
                        <Share2 size={18} />
                    </button>
                </div>
                {/* Report button */}
                <button className="p-1.5 rounded-full hover:bg-gray-100 transition-colors text-gray-500">
                    <Flag size={18} />
                </button>
            </div>
        </div>
    );
};
const SimplifiedPostCardPreview = () => {
    return (
        <div className="p-4 bg-gray-100 min-h-screen flex justify-center items-center">
            <SimplifiedPostCard />
        </div>
    );
};
export default SimplifiedPostCardPreview;