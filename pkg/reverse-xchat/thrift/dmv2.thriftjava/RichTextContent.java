package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import com.google.android.libraries.places.api.model.PlaceTypes;
import java.io.IOException;
import kotlin.Metadata;
import kotlin.NoWhenBranchMatchedException;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.DefaultConstructorMarker;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u0000<\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\f\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\b6\u0018\u0000 \t2\u00020\u0001:\n\n\u000b\f\r\u000e\u000f\u0010\u0011\u0012\tB\t\b\u0004¢\u0006\u0004\b\u0002\u0010\u0003J\u0017\u0010\u0007\u001a\u00020\u00062\u0006\u0010\u0005\u001a\u00020\u0004H\u0016¢\u0006\u0004\b\u0007\u0010\b\u0082\u0001\b\u0013\u0014\u0015\u0016\u0017\u0018\u0019\u001a¨\u0006\u001b"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/RichTextContent;", "Lcom/bendb/thrifty/a;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "Companion", "Hashtag", "Cashtag", "Mention", "Url", "Email", "Address", "PhoneNumber", "Unknown", "RichTextContentAdapter", "Lcom/x/dmv2/thriftjava/RichTextContent$Address;", "Lcom/x/dmv2/thriftjava/RichTextContent$Cashtag;", "Lcom/x/dmv2/thriftjava/RichTextContent$Email;", "Lcom/x/dmv2/thriftjava/RichTextContent$Hashtag;", "Lcom/x/dmv2/thriftjava/RichTextContent$Mention;", "Lcom/x/dmv2/thriftjava/RichTextContent$PhoneNumber;", "Lcom/x/dmv2/thriftjava/RichTextContent$Unknown;", "Lcom/x/dmv2/thriftjava/RichTextContent$Url;", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public abstract class RichTextContent implements InterfaceC11261a {

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new RichTextContentAdapter();

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/RichTextContent$Address;", "Lcom/x/dmv2/thriftjava/RichTextContent;", "value", "Lcom/x/dmv2/thriftjava/AddressRichTextContent;", "<init>", "(Lcom/x/dmv2/thriftjava/AddressRichTextContent;)V", "getValue", "()Lcom/x/dmv2/thriftjava/AddressRichTextContent;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Address extends RichTextContent {

        @InterfaceC88464a
        private final AddressRichTextContent value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public Address(@InterfaceC88464a AddressRichTextContent value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ Address copy$default(Address address, AddressRichTextContent addressRichTextContent, int i, Object obj) {
            if ((i & 1) != 0) {
                addressRichTextContent = address.value;
            }
            return address.copy(addressRichTextContent);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final AddressRichTextContent getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final Address copy(@InterfaceC88464a AddressRichTextContent value) {
            Intrinsics.m65272h(value, "value");
            return new Address(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof Address) && Intrinsics.m65267c(this.value, ((Address) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final AddressRichTextContent m76811getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "RichTextContent(address=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/RichTextContent$Cashtag;", "Lcom/x/dmv2/thriftjava/RichTextContent;", "value", "Lcom/x/dmv2/thriftjava/CashtagRichTextContent;", "<init>", "(Lcom/x/dmv2/thriftjava/CashtagRichTextContent;)V", "getValue", "()Lcom/x/dmv2/thriftjava/CashtagRichTextContent;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Cashtag extends RichTextContent {

        @InterfaceC88464a
        private final CashtagRichTextContent value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public Cashtag(@InterfaceC88464a CashtagRichTextContent value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ Cashtag copy$default(Cashtag cashtag, CashtagRichTextContent cashtagRichTextContent, int i, Object obj) {
            if ((i & 1) != 0) {
                cashtagRichTextContent = cashtag.value;
            }
            return cashtag.copy(cashtagRichTextContent);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final CashtagRichTextContent getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final Cashtag copy(@InterfaceC88464a CashtagRichTextContent value) {
            Intrinsics.m65272h(value, "value");
            return new Cashtag(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof Cashtag) && Intrinsics.m65267c(this.value, ((Cashtag) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final CashtagRichTextContent m76812getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "RichTextContent(cashtag=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/RichTextContent$Email;", "Lcom/x/dmv2/thriftjava/RichTextContent;", "value", "Lcom/x/dmv2/thriftjava/EmailRichTextContent;", "<init>", "(Lcom/x/dmv2/thriftjava/EmailRichTextContent;)V", "getValue", "()Lcom/x/dmv2/thriftjava/EmailRichTextContent;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Email extends RichTextContent {

        @InterfaceC88464a
        private final EmailRichTextContent value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public Email(@InterfaceC88464a EmailRichTextContent value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ Email copy$default(Email email, EmailRichTextContent emailRichTextContent, int i, Object obj) {
            if ((i & 1) != 0) {
                emailRichTextContent = email.value;
            }
            return email.copy(emailRichTextContent);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final EmailRichTextContent getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final Email copy(@InterfaceC88464a EmailRichTextContent value) {
            Intrinsics.m65272h(value, "value");
            return new Email(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof Email) && Intrinsics.m65267c(this.value, ((Email) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final EmailRichTextContent m76813getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "RichTextContent(email=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/RichTextContent$Hashtag;", "Lcom/x/dmv2/thriftjava/RichTextContent;", "value", "Lcom/x/dmv2/thriftjava/HashtagRichTextContent;", "<init>", "(Lcom/x/dmv2/thriftjava/HashtagRichTextContent;)V", "getValue", "()Lcom/x/dmv2/thriftjava/HashtagRichTextContent;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Hashtag extends RichTextContent {

        @InterfaceC88464a
        private final HashtagRichTextContent value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public Hashtag(@InterfaceC88464a HashtagRichTextContent value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ Hashtag copy$default(Hashtag hashtag, HashtagRichTextContent hashtagRichTextContent, int i, Object obj) {
            if ((i & 1) != 0) {
                hashtagRichTextContent = hashtag.value;
            }
            return hashtag.copy(hashtagRichTextContent);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final HashtagRichTextContent getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final Hashtag copy(@InterfaceC88464a HashtagRichTextContent value) {
            Intrinsics.m65272h(value, "value");
            return new Hashtag(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof Hashtag) && Intrinsics.m65267c(this.value, ((Hashtag) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final HashtagRichTextContent m76814getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "RichTextContent(hashtag=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/RichTextContent$Mention;", "Lcom/x/dmv2/thriftjava/RichTextContent;", "value", "Lcom/x/dmv2/thriftjava/MentionRichTextContent;", "<init>", "(Lcom/x/dmv2/thriftjava/MentionRichTextContent;)V", "getValue", "()Lcom/x/dmv2/thriftjava/MentionRichTextContent;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Mention extends RichTextContent {

        @InterfaceC88464a
        private final MentionRichTextContent value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public Mention(@InterfaceC88464a MentionRichTextContent value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ Mention copy$default(Mention mention, MentionRichTextContent mentionRichTextContent, int i, Object obj) {
            if ((i & 1) != 0) {
                mentionRichTextContent = mention.value;
            }
            return mention.copy(mentionRichTextContent);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final MentionRichTextContent getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final Mention copy(@InterfaceC88464a MentionRichTextContent value) {
            Intrinsics.m65272h(value, "value");
            return new Mention(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof Mention) && Intrinsics.m65267c(this.value, ((Mention) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final MentionRichTextContent m76815getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "RichTextContent(mention=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/RichTextContent$PhoneNumber;", "Lcom/x/dmv2/thriftjava/RichTextContent;", "value", "Lcom/x/dmv2/thriftjava/PhoneNumberRichTextContent;", "<init>", "(Lcom/x/dmv2/thriftjava/PhoneNumberRichTextContent;)V", "getValue", "()Lcom/x/dmv2/thriftjava/PhoneNumberRichTextContent;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class PhoneNumber extends RichTextContent {

        @InterfaceC88464a
        private final PhoneNumberRichTextContent value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public PhoneNumber(@InterfaceC88464a PhoneNumberRichTextContent value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ PhoneNumber copy$default(PhoneNumber phoneNumber, PhoneNumberRichTextContent phoneNumberRichTextContent, int i, Object obj) {
            if ((i & 1) != 0) {
                phoneNumberRichTextContent = phoneNumber.value;
            }
            return phoneNumber.copy(phoneNumberRichTextContent);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final PhoneNumberRichTextContent getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final PhoneNumber copy(@InterfaceC88464a PhoneNumberRichTextContent value) {
            Intrinsics.m65272h(value, "value");
            return new PhoneNumber(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof PhoneNumber) && Intrinsics.m65267c(this.value, ((PhoneNumber) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final PhoneNumberRichTextContent m76816getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "RichTextContent(phoneNumber=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/RichTextContent$RichTextContentAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/RichTextContent;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/RichTextContent;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/RichTextContent;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class RichTextContentAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public RichTextContent m83779read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            RichTextContent hashtag;
            Intrinsics.m65272h(protocol, "protocol");
            RichTextContent richTextContent = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    if (richTextContent != null) {
                        return richTextContent;
                    }
                    throw new IllegalStateException("unreadable");
                }
                switch (c11265cMo14127V2.f38393b) {
                    case 1:
                        if (b == 12) {
                            hashtag = new Hashtag((HashtagRichTextContent) HashtagRichTextContent.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 2:
                        if (b == 12) {
                            hashtag = new Cashtag((CashtagRichTextContent) CashtagRichTextContent.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 3:
                        if (b == 12) {
                            hashtag = new Mention((MentionRichTextContent) MentionRichTextContent.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 4:
                        if (b == 12) {
                            hashtag = new Url((UrlRichTextContent) UrlRichTextContent.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 5:
                        if (b == 12) {
                            hashtag = new Email((EmailRichTextContent) EmailRichTextContent.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 6:
                        if (b == 12) {
                            hashtag = new Address((AddressRichTextContent) AddressRichTextContent.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 7:
                        if (b == 12) {
                            hashtag = new PhoneNumber((PhoneNumberRichTextContent) PhoneNumberRichTextContent.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    default:
                        richTextContent = Unknown.INSTANCE;
                        C11272a.m14141a(protocol, b);
                        continue;
                }
                richTextContent = hashtag;
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a RichTextContent struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("RichTextContent");
            if (struct instanceof Hashtag) {
                protocol.mo14136v3("hashtag", 1, (byte) 12);
                HashtagRichTextContent.ADAPTER.write(protocol, ((Hashtag) struct).m76814getValue());
            } else if (struct instanceof Cashtag) {
                protocol.mo14136v3("cashtag", 2, (byte) 12);
                CashtagRichTextContent.ADAPTER.write(protocol, ((Cashtag) struct).m76812getValue());
            } else if (struct instanceof Mention) {
                protocol.mo14136v3("mention", 3, (byte) 12);
                MentionRichTextContent.ADAPTER.write(protocol, ((Mention) struct).m76815getValue());
            } else if (struct instanceof Url) {
                protocol.mo14136v3("url", 4, (byte) 12);
                UrlRichTextContent.ADAPTER.write(protocol, ((Url) struct).m76817getValue());
            } else if (struct instanceof Email) {
                protocol.mo14136v3("email", 5, (byte) 12);
                EmailRichTextContent.ADAPTER.write(protocol, ((Email) struct).m76813getValue());
            } else if (struct instanceof Address) {
                protocol.mo14136v3(PlaceTypes.ADDRESS, 6, (byte) 12);
                AddressRichTextContent.ADAPTER.write(protocol, ((Address) struct).m76811getValue());
            } else if (struct instanceof PhoneNumber) {
                protocol.mo14136v3("phoneNumber", 7, (byte) 12);
                PhoneNumberRichTextContent.ADAPTER.write(protocol, ((PhoneNumber) struct).m76816getValue());
            } else if (!(struct instanceof Unknown)) {
                throw new NoWhenBranchMatchedException();
            }
            protocol.mo14134i0();
        }
    }

    @Metadata(m64929d1 = {"\u0000$\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\n\u0002\u0010\u000e\n\u0000\bÆ\n\u0018\u00002\u00020\u0001B\t\b\u0002¢\u0006\u0004\b\u0002\u0010\u0003J\u0013\u0010\u0004\u001a\u00020\u00052\b\u0010\u0006\u001a\u0004\u0018\u00010\u0007HÖ\u0003J\t\u0010\b\u001a\u00020\tHÖ\u0001J\t\u0010\n\u001a\u00020\u000bHÖ\u0001¨\u0006\f"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/RichTextContent$Unknown;", "Lcom/x/dmv2/thriftjava/RichTextContent;", "<init>", "()V", "equals", "", "other", "", "hashCode", "", "toString", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Unknown extends RichTextContent {

        @InterfaceC88464a
        public static final Unknown INSTANCE = new Unknown();

        private Unknown() {
            super(null);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            return this == other || (other instanceof Unknown);
        }

        public int hashCode() {
            return -48506689;
        }

        @InterfaceC88464a
        public String toString() {
            return "Unknown";
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/RichTextContent$Url;", "Lcom/x/dmv2/thriftjava/RichTextContent;", "value", "Lcom/x/dmv2/thriftjava/UrlRichTextContent;", "<init>", "(Lcom/x/dmv2/thriftjava/UrlRichTextContent;)V", "getValue", "()Lcom/x/dmv2/thriftjava/UrlRichTextContent;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Url extends RichTextContent {

        @InterfaceC88464a
        private final UrlRichTextContent value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public Url(@InterfaceC88464a UrlRichTextContent value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ Url copy$default(Url url, UrlRichTextContent urlRichTextContent, int i, Object obj) {
            if ((i & 1) != 0) {
                urlRichTextContent = url.value;
            }
            return url.copy(urlRichTextContent);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final UrlRichTextContent getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final Url copy(@InterfaceC88464a UrlRichTextContent value) {
            Intrinsics.m65272h(value, "value");
            return new Url(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof Url) && Intrinsics.m65267c(this.value, ((Url) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final UrlRichTextContent m76817getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "RichTextContent(url=" + this.value + Separators.RPAREN;
        }
    }

    public /* synthetic */ RichTextContent(DefaultConstructorMarker defaultConstructorMarker) {
        this();
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }

    private RichTextContent() {
    }
}
