/*eslint-disable block-scoped-var, id-length, no-control-regex, no-magic-numbers, no-prototype-builtins, no-redeclare, no-shadow, no-var, sort-vars*/
import * as $protobuf from "protobufjs/minimal";
// Common aliases
const $Reader = $protobuf.Reader, $Writer = $protobuf.Writer, $util = $protobuf.util;
// Exported root namespace
const $root = $protobuf.roots["default"] || ($protobuf.roots["default"] = {});
export const mes = $root.mes = (() => {
    /**
     * Namespace mes.
     * @exports mes
     * @namespace
     */
    const mes = {};
    mes.SendMessage = (function() {
        /**
         * Properties of a SendMessage.
         * @memberof mes
         * @interface ISendMessage
         * @property {string|null} [id] SendMessage id
         * @property {string|null} [conversationId] SendMessage conversationId
         * @property {string|null} [senderId] SendMessage senderId
         * @property {string|null} [recipientId] SendMessage recipientId
         * @property {string|null} [itemId] SendMessage itemId
         * @property {string|null} [body] SendMessage body
         * @property {boolean|null} [isRead] SendMessage isRead
         */
        /**
         * Constructs a new SendMessage.
         * @memberof mes
         * @classdesc Represents a SendMessage.
         * @implements ISendMessage
         * @constructor
         * @param {mes.ISendMessage=} [properties] Properties to set
         */
        function SendMessage(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }
        /**
         * SendMessage id.
         * @member {string} id
         * @memberof mes.SendMessage
         * @instance
         */
        SendMessage.prototype.id = "";
        /**
         * SendMessage conversationId.
         * @member {string} conversationId
         * @memberof mes.SendMessage
         * @instance
         */
        SendMessage.prototype.conversationId = "";
        /**
         * SendMessage senderId.
         * @member {string} senderId
         * @memberof mes.SendMessage
         * @instance
         */
        SendMessage.prototype.senderId = "";
        /**
         * SendMessage recipientId.
         * @member {string} recipientId
         * @memberof mes.SendMessage
         * @instance
         */
        SendMessage.prototype.recipientId = "";
        /**
         * SendMessage itemId.
         * @member {string} itemId
         * @memberof mes.SendMessage
         * @instance
         */
        SendMessage.prototype.itemId = "";
        /**
         * SendMessage body.
         * @member {string} body
         * @memberof mes.SendMessage
         * @instance
         */
        SendMessage.prototype.body = "";
        /**
         * SendMessage isRead.
         * @member {boolean} isRead
         * @memberof mes.SendMessage
         * @instance
         */
        SendMessage.prototype.isRead = false;
        /**
         * Creates a new SendMessage instance using the specified properties.
         * @function create
         * @memberof mes.SendMessage
         * @static
         * @param {mes.ISendMessage=} [properties] Properties to set
         * @returns {mes.SendMessage} SendMessage instance
         */
        SendMessage.create = function create(properties) {
            return new SendMessage(properties);
        };
        /**
         * Encodes the specified SendMessage message. Does not implicitly {@link mes.SendMessage.verify|verify} messages.
         * @function encode
         * @memberof mes.SendMessage
         * @static
         * @param {mes.ISendMessage} message SendMessage message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        SendMessage.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.id != null && Object.hasOwnProperty.call(message, "id"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
            if (message.conversationId != null && Object.hasOwnProperty.call(message, "conversationId"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.conversationId);
            if (message.senderId != null && Object.hasOwnProperty.call(message, "senderId"))
                writer.uint32(/* id 3, wireType 2 =*/26).string(message.senderId);
            if (message.recipientId != null && Object.hasOwnProperty.call(message, "recipientId"))
                writer.uint32(/* id 4, wireType 2 =*/34).string(message.recipientId);
            if (message.itemId != null && Object.hasOwnProperty.call(message, "itemId"))
                writer.uint32(/* id 5, wireType 2 =*/42).string(message.itemId);
            if (message.body != null && Object.hasOwnProperty.call(message, "body"))
                writer.uint32(/* id 6, wireType 2 =*/50).string(message.body);
            if (message.isRead != null && Object.hasOwnProperty.call(message, "isRead"))
                writer.uint32(/* id 7, wireType 0 =*/56).bool(message.isRead);
            return writer;
        };
        /**
         * Encodes the specified SendMessage message, length delimited. Does not implicitly {@link mes.SendMessage.verify|verify} messages.
         * @function encodeDelimited
         * @memberof mes.SendMessage
         * @static
         * @param {mes.ISendMessage} message SendMessage message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        SendMessage.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };
        /**
         * Decodes a SendMessage message from the specified reader or buffer.
         * @function decode
         * @memberof mes.SendMessage
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {mes.SendMessage} SendMessage
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        SendMessage.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.mes.SendMessage();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        message.id = reader.string();
                        break;
                    }
                case 2: {
                        message.conversationId = reader.string();
                        break;
                    }
                case 3: {
                        message.senderId = reader.string();
                        break;
                    }
                case 4: {
                        message.recipientId = reader.string();
                        break;
                    }
                case 5: {
                        message.itemId = reader.string();
                        break;
                    }
                case 6: {
                        message.body = reader.string();
                        break;
                    }
                case 7: {
                        message.isRead = reader.bool();
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };
        /**
         * Decodes a SendMessage message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof mes.SendMessage
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {mes.SendMessage} SendMessage
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        SendMessage.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };
        /**
         * Verifies a SendMessage message.
         * @function verify
         * @memberof mes.SendMessage
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        SendMessage.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.id != null && message.hasOwnProperty("id"))
                if (!$util.isString(message.id))
                    return "id: string expected";
            if (message.conversationId != null && message.hasOwnProperty("conversationId"))
                if (!$util.isString(message.conversationId))
                    return "conversationId: string expected";
            if (message.senderId != null && message.hasOwnProperty("senderId"))
                if (!$util.isString(message.senderId))
                    return "senderId: string expected";
            if (message.recipientId != null && message.hasOwnProperty("recipientId"))
                if (!$util.isString(message.recipientId))
                    return "recipientId: string expected";
            if (message.itemId != null && message.hasOwnProperty("itemId"))
                if (!$util.isString(message.itemId))
                    return "itemId: string expected";
            if (message.body != null && message.hasOwnProperty("body"))
                if (!$util.isString(message.body))
                    return "body: string expected";
            if (message.isRead != null && message.hasOwnProperty("isRead"))
                if (typeof message.isRead !== "boolean")
                    return "isRead: boolean expected";
            return null;
        };
        /**
         * Creates a SendMessage message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof mes.SendMessage
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {mes.SendMessage} SendMessage
         */
        SendMessage.fromObject = function fromObject(object) {
            if (object instanceof $root.mes.SendMessage)
                return object;
            let message = new $root.mes.SendMessage();
            if (object.id != null)
                message.id = String(object.id);
            if (object.conversationId != null)
                message.conversationId = String(object.conversationId);
            if (object.senderId != null)
                message.senderId = String(object.senderId);
            if (object.recipientId != null)
                message.recipientId = String(object.recipientId);
            if (object.itemId != null)
                message.itemId = String(object.itemId);
            if (object.body != null)
                message.body = String(object.body);
            if (object.isRead != null)
                message.isRead = Boolean(object.isRead);
            return message;
        };
        /**
         * Creates a plain object from a SendMessage message. Also converts values to other types if specified.
         * @function toObject
         * @memberof mes.SendMessage
         * @static
         * @param {mes.SendMessage} message SendMessage
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        SendMessage.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.id = "";
                object.conversationId = "";
                object.senderId = "";
                object.recipientId = "";
                object.itemId = "";
                object.body = "";
                object.isRead = false;
            }
            if (message.id != null && message.hasOwnProperty("id"))
                object.id = message.id;
            if (message.conversationId != null && message.hasOwnProperty("conversationId"))
                object.conversationId = message.conversationId;
            if (message.senderId != null && message.hasOwnProperty("senderId"))
                object.senderId = message.senderId;
            if (message.recipientId != null && message.hasOwnProperty("recipientId"))
                object.recipientId = message.recipientId;
            if (message.itemId != null && message.hasOwnProperty("itemId"))
                object.itemId = message.itemId;
            if (message.body != null && message.hasOwnProperty("body"))
                object.body = message.body;
            if (message.isRead != null && message.hasOwnProperty("isRead"))
                object.isRead = message.isRead;
            return object;
        };
        /**
         * Converts this SendMessage to JSON.
         * @function toJSON
         * @memberof mes.SendMessage
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        SendMessage.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };
        /**
         * Gets the default type url for SendMessage
         * @function getTypeUrl
         * @memberof mes.SendMessage
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        SendMessage.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/mes.SendMessage";
        };
        return SendMessage;
    })();
    return mes;
})();
export { $root as default };
