package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import java.io.IOException;
import kotlin.Metadata;
import kotlin.NoWhenBranchMatchedException;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.DefaultConstructorMarker;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u00004\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\n\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\b6\u0018\u0000 \t2\u00020\u0001:\b\n\u000b\f\r\u000e\u000f\u0010\tB\t\b\u0004¢\u0006\u0004\b\u0002\u0010\u0003J\u0017\u0010\u0007\u001a\u00020\u00062\u0006\u0010\u0005\u001a\u00020\u0004H\u0016¢\u0006\u0004\b\u0007\u0010\b\u0082\u0001\u0006\u0011\u0012\u0013\u0014\u0015\u0016¨\u0006\u0017"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageAttachment;", "Lcom/bendb/thrifty/a;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "Companion", "Media", "Post", "Url", "UnifiedCard", "Money", "Unknown", "MessageAttachmentAdapter", "Lcom/x/dmv2/thriftjava/MessageAttachment$Media;", "Lcom/x/dmv2/thriftjava/MessageAttachment$Money;", "Lcom/x/dmv2/thriftjava/MessageAttachment$Post;", "Lcom/x/dmv2/thriftjava/MessageAttachment$UnifiedCard;", "Lcom/x/dmv2/thriftjava/MessageAttachment$Unknown;", "Lcom/x/dmv2/thriftjava/MessageAttachment$Url;", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public abstract class MessageAttachment implements InterfaceC11261a {

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new MessageAttachmentAdapter();

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageAttachment$Media;", "Lcom/x/dmv2/thriftjava/MessageAttachment;", "value", "Lcom/x/dmv2/thriftjava/MediaAttachment;", "<init>", "(Lcom/x/dmv2/thriftjava/MediaAttachment;)V", "getValue", "()Lcom/x/dmv2/thriftjava/MediaAttachment;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Media extends MessageAttachment {

        @InterfaceC88464a
        private final MediaAttachment value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public Media(@InterfaceC88464a MediaAttachment value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ Media copy$default(Media media, MediaAttachment mediaAttachment, int i, Object obj) {
            if ((i & 1) != 0) {
                mediaAttachment = media.value;
            }
            return media.copy(mediaAttachment);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final MediaAttachment getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final Media copy(@InterfaceC88464a MediaAttachment value) {
            Intrinsics.m65272h(value, "value");
            return new Media(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof Media) && Intrinsics.m65267c(this.value, ((Media) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final MediaAttachment m76763getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageAttachment(media=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageAttachment$MessageAttachmentAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/MessageAttachment;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/MessageAttachment;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/MessageAttachment;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class MessageAttachmentAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public MessageAttachment m85949read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            MessageAttachment money;
            Intrinsics.m65272h(protocol, "protocol");
            MessageAttachment messageAttachment = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    break;
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        if (s != 3) {
                            if (s != 4) {
                                if (s != 5) {
                                    messageAttachment = Unknown.INSTANCE;
                                    C11272a.m14141a(protocol, b);
                                } else if (b == 12) {
                                    money = new Money((MoneyAttachment) MoneyAttachment.ADAPTER.read(protocol));
                                    messageAttachment = money;
                                } else {
                                    C11272a.m14141a(protocol, b);
                                }
                            } else if (b == 12) {
                                money = new UnifiedCard((UnifiedCardAttachment) UnifiedCardAttachment.ADAPTER.read(protocol));
                                messageAttachment = money;
                            } else {
                                C11272a.m14141a(protocol, b);
                            }
                        } else if (b == 12) {
                            money = new Url((UrlAttachment) UrlAttachment.ADAPTER.read(protocol));
                            messageAttachment = money;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    } else if (b == 12) {
                        money = new Post((PostAttachment) PostAttachment.ADAPTER.read(protocol));
                        messageAttachment = money;
                    } else {
                        C11272a.m14141a(protocol, b);
                    }
                } else if (b == 12) {
                    money = new Media((MediaAttachment) MediaAttachment.ADAPTER.read(protocol));
                    messageAttachment = money;
                } else {
                    C11272a.m14141a(protocol, b);
                }
            }
            if (messageAttachment != null) {
                return messageAttachment;
            }
            throw new IllegalStateException("unreadable");
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a MessageAttachment struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("MessageAttachment");
            if (struct instanceof Media) {
                protocol.mo14136v3("media", 1, (byte) 12);
                MediaAttachment.ADAPTER.write(protocol, ((Media) struct).m76763getValue());
            } else if (struct instanceof Post) {
                protocol.mo14136v3("post", 2, (byte) 12);
                PostAttachment.ADAPTER.write(protocol, ((Post) struct).m76765getValue());
            } else if (struct instanceof Url) {
                protocol.mo14136v3("url", 3, (byte) 12);
                UrlAttachment.ADAPTER.write(protocol, ((Url) struct).m76767getValue());
            } else if (struct instanceof UnifiedCard) {
                protocol.mo14136v3("unified_card", 4, (byte) 12);
                UnifiedCardAttachment.ADAPTER.write(protocol, ((UnifiedCard) struct).m76766getValue());
            } else if (struct instanceof Money) {
                protocol.mo14136v3("money", 5, (byte) 12);
                MoneyAttachment.ADAPTER.write(protocol, ((Money) struct).m76764getValue());
            } else if (!(struct instanceof Unknown)) {
                throw new NoWhenBranchMatchedException();
            }
            protocol.mo14134i0();
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageAttachment$Money;", "Lcom/x/dmv2/thriftjava/MessageAttachment;", "value", "Lcom/x/dmv2/thriftjava/MoneyAttachment;", "<init>", "(Lcom/x/dmv2/thriftjava/MoneyAttachment;)V", "getValue", "()Lcom/x/dmv2/thriftjava/MoneyAttachment;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Money extends MessageAttachment {

        @InterfaceC88464a
        private final MoneyAttachment value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public Money(@InterfaceC88464a MoneyAttachment value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ Money copy$default(Money money, MoneyAttachment moneyAttachment, int i, Object obj) {
            if ((i & 1) != 0) {
                moneyAttachment = money.value;
            }
            return money.copy(moneyAttachment);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final MoneyAttachment getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final Money copy(@InterfaceC88464a MoneyAttachment value) {
            Intrinsics.m65272h(value, "value");
            return new Money(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof Money) && Intrinsics.m65267c(this.value, ((Money) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final MoneyAttachment m76764getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageAttachment(money=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageAttachment$Post;", "Lcom/x/dmv2/thriftjava/MessageAttachment;", "value", "Lcom/x/dmv2/thriftjava/PostAttachment;", "<init>", "(Lcom/x/dmv2/thriftjava/PostAttachment;)V", "getValue", "()Lcom/x/dmv2/thriftjava/PostAttachment;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Post extends MessageAttachment {

        @InterfaceC88464a
        private final PostAttachment value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public Post(@InterfaceC88464a PostAttachment value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ Post copy$default(Post post, PostAttachment postAttachment, int i, Object obj) {
            if ((i & 1) != 0) {
                postAttachment = post.value;
            }
            return post.copy(postAttachment);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final PostAttachment getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final Post copy(@InterfaceC88464a PostAttachment value) {
            Intrinsics.m65272h(value, "value");
            return new Post(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof Post) && Intrinsics.m65267c(this.value, ((Post) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final PostAttachment m76765getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageAttachment(post=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageAttachment$UnifiedCard;", "Lcom/x/dmv2/thriftjava/MessageAttachment;", "value", "Lcom/x/dmv2/thriftjava/UnifiedCardAttachment;", "<init>", "(Lcom/x/dmv2/thriftjava/UnifiedCardAttachment;)V", "getValue", "()Lcom/x/dmv2/thriftjava/UnifiedCardAttachment;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class UnifiedCard extends MessageAttachment {

        @InterfaceC88464a
        private final UnifiedCardAttachment value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public UnifiedCard(@InterfaceC88464a UnifiedCardAttachment value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ UnifiedCard copy$default(UnifiedCard unifiedCard, UnifiedCardAttachment unifiedCardAttachment, int i, Object obj) {
            if ((i & 1) != 0) {
                unifiedCardAttachment = unifiedCard.value;
            }
            return unifiedCard.copy(unifiedCardAttachment);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final UnifiedCardAttachment getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final UnifiedCard copy(@InterfaceC88464a UnifiedCardAttachment value) {
            Intrinsics.m65272h(value, "value");
            return new UnifiedCard(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof UnifiedCard) && Intrinsics.m65267c(this.value, ((UnifiedCard) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final UnifiedCardAttachment m76766getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageAttachment(unified_card=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000$\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\n\u0002\u0010\u000e\n\u0000\bÆ\n\u0018\u00002\u00020\u0001B\t\b\u0002¢\u0006\u0004\b\u0002\u0010\u0003J\u0013\u0010\u0004\u001a\u00020\u00052\b\u0010\u0006\u001a\u0004\u0018\u00010\u0007HÖ\u0003J\t\u0010\b\u001a\u00020\tHÖ\u0001J\t\u0010\n\u001a\u00020\u000bHÖ\u0001¨\u0006\f"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageAttachment$Unknown;", "Lcom/x/dmv2/thriftjava/MessageAttachment;", "<init>", "()V", "equals", "", "other", "", "hashCode", "", "toString", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Unknown extends MessageAttachment {

        @InterfaceC88464a
        public static final Unknown INSTANCE = new Unknown();

        private Unknown() {
            super(null);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            return this == other || (other instanceof Unknown);
        }

        public int hashCode() {
            return -419406215;
        }

        @InterfaceC88464a
        public String toString() {
            return "Unknown";
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageAttachment$Url;", "Lcom/x/dmv2/thriftjava/MessageAttachment;", "value", "Lcom/x/dmv2/thriftjava/UrlAttachment;", "<init>", "(Lcom/x/dmv2/thriftjava/UrlAttachment;)V", "getValue", "()Lcom/x/dmv2/thriftjava/UrlAttachment;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Url extends MessageAttachment {

        @InterfaceC88464a
        private final UrlAttachment value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public Url(@InterfaceC88464a UrlAttachment value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ Url copy$default(Url url, UrlAttachment urlAttachment, int i, Object obj) {
            if ((i & 1) != 0) {
                urlAttachment = url.value;
            }
            return url.copy(urlAttachment);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final UrlAttachment getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final Url copy(@InterfaceC88464a UrlAttachment value) {
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
        public final UrlAttachment m76767getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageAttachment(url=" + this.value + Separators.RPAREN;
        }
    }

    public /* synthetic */ MessageAttachment(DefaultConstructorMarker defaultConstructorMarker) {
        this();
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }

    private MessageAttachment() {
    }
}